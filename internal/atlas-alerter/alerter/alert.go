package alerter

import (
	"bytes"
	"context"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/template"
	"google.golang.org/protobuf/proto"
)

// ErrUnknownAlarm is returned when sending an alarm of an unknown type.
const ErrUnknownAlarm consterr.Error = "unknown alarm type"

// alertMessages receives event messages, alerts based on alarm processing, and
// stores the results.
func (ale *Alerter) alertMessages() {
	alog.Info("alertMessages starting processor")
	ctx := context.Background()

	var processCount int
	for msg := range ale.eOutSub.C() {
		// Retrieve published message.
		metric.Incr("received", nil)
		eOut := &message.EventerOut{}
		err := proto.Unmarshal(msg.Payload(), eOut)
		if err != nil || eOut.Point == nil || eOut.Device == nil ||
			eOut.Rule == nil {
			msg.Ack()

			if !bytes.Equal([]byte{queue.Prime}, msg.Payload()) {
				metric.Incr("error", map[string]string{"func": "unmarshal"})
				alog.Errorf("alertMessages proto.Unmarshal eOut, err: %+v, %v",
					eOut, err)
			}

			continue
		}

		// Set up logging fields.
		logger := alog.
			WithField("traceID", eOut.Point.TraceId).
			WithField("orgID", eOut.Device.OrgId).
			WithField("uniqID", eOut.Point.UniqId).
			WithField("devID", eOut.Device.Id)

		// Retrieve org.
		dCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		org, err := ale.orgDAO.Read(dCtx, eOut.Device.OrgId)
		cancel()
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "read"})
			logger.Errorf("alertMessages ale.orgDAO.Read: %v", err)

			continue
		}
		logger.Debugf("alertMessages org: %+v", org)

		// Retrieve alarms by rule ID. Alarms may be disabled.
		dCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
		alarms, _, err := ale.alarmDAO.List(dCtx, eOut.Device.OrgId,
			time.Time{}, "", 0, eOut.Rule.Id)
		cancel()
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "list"})
			logger.Errorf("alertMessages ale.alarmDAO.List: %v", err)

			continue
		}
		logger.Debugf("alertMessages alarms: %+v", alarms)

		// Validate, retrieve users, process and send alerts, and store results.
		for _, a := range alarms {
			ale.evalAlarms(alog.NewContext(
				ctx, &alog.CtxLogger{Logger: logger}), eOut, org, a)
		}

		msg.Ack()
		metric.Incr("processed", nil)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("alertMessages processed %v messages", processCount)
		}
	}
}

// evalAlarms validates alarms, retrieves users, processes and sends alerts, and
// stores results. Unconditionally acknowledge a message after processing, as
// there are no guarantees of alarms or users being assigned to an event.
func (ale *Alerter) evalAlarms(
	ctx context.Context, eOut *message.EventerOut, org *api.Org, a *api.Alarm,
) {
	logger := alog.FromContext(ctx)

	// Validate alarm.
	if a.Status != api.Status_ACTIVE {
		return
	}

	// Retrieve users. Only active users with matching tags will be returned.
	dCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	users, err := ale.userDAO.ListByTags(dCtx, eOut.Device.OrgId, a.UserTags)
	cancel()
	if err != nil {
		metric.Incr("error", map[string]string{"func": "listbytags"})
		logger.Errorf("alertMessages ale.userDAO.ListByTags: %v", err)

		return
	}
	if len(users) == 0 {
		return
	}

	// Generate alert subject and body.
	subj, err := template.Generate(eOut.Point, eOut.Rule, eOut.Device,
		a.SubjectTemplate)
	if err != nil {
		metric.Incr("error", map[string]string{"func": "gensubject"})
		logger.Errorf("alertMessages subject template.Generate: %v", err)

		return
	}

	body, err := template.Generate(eOut.Point, eOut.Rule, eOut.Device,
		a.BodyTemplate)
	if err != nil {
		metric.Incr("error", map[string]string{"func": "genbody"})
		logger.Errorf("alertMessages body template.Generate: %v", err)

		return
	}

	// Process alerts.
	for _, user := range users {
		// Check cache for existing repeat interval.
		key := repeatKey(eOut.Device.OrgId, eOut.Device.Id, a.Id, user.Id)
		cCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		ok, err := ale.cache.SetIfNotExistTTL(cCtx, key, 1,
			time.Duration(a.RepeatInterval)*time.Minute)
		cancel()
		if err != nil {
			metric.Incr("error", map[string]string{"func": "setifnotexist"})
			logger.Errorf("alertMessages ale.cache.SetIfNotExistTTL: %v", err)

			continue
		}
		if !ok {
			metric.Incr("repeat", nil)

			continue
		}

		// Send alert.
		nCtx, cancel := context.WithTimeout(ctx, time.Minute)
		switch a.Type {
		case api.AlarmType_APP:
			err = ale.notify.App(nCtx, user.AppKey, subj, body)
		case api.AlarmType_SMS:
			err = ale.notify.SMS(nCtx, user.Phone, subj, body)
		case api.AlarmType_EMAIL:
			err = ale.notify.Email(nCtx, org.DisplayName, org.Email, user.Email,
				subj, body)
		case api.AlarmType_ALARM_TYPE_UNSPECIFIED:
			fallthrough
		default:
			err = ErrUnknownAlarm
		}
		cancel()

		alert := &api.Alert{
			OrgId:   eOut.Device.OrgId,
			UniqId:  eOut.Device.UniqId,
			AlarmId: a.Id,
			UserId:  user.Id,
			TraceId: eOut.Point.TraceId,
		}

		if err != nil {
			alert.Status = api.AlertStatus_ERROR
			alert.Error = err.Error()
			metric.Incr("error", map[string]string{"func": "notify"})
			logger.Errorf("alertMessages ale.notify a, err: %+v, %v", a,
				err.Error())
		} else {
			alert.Status = api.AlertStatus_SENT
			metric.Incr("sent", map[string]string{
				"type": a.Type.String(),
			})
			logger.Debugf("alertMessages sent user, msg: %+v, %v", user,
				subj+" - "+body)
		}

		// Store alert.
		dCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
		err = ale.aleDAO.Create(dCtx, alert)
		cancel()
		if err != nil {
			metric.Incr("error", map[string]string{"func": "create"})
			logger.Errorf("alertMessages ale.aleDAO.Create: %v", err)

			continue
		}

		metric.Incr("created", nil)
		logger.Debugf("alertMessages created: %+v", alert)
	}
}
