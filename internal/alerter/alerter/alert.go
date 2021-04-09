package alerter

import (
	"bytes"
	"context"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alarm"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"google.golang.org/protobuf/proto"
)

const ErrUnknownAlarm consterr.Error = "unknown alarm type"

// alertMessages receives event messages, alerts based on alarm processing, and
// stores the results.
func (ale *Alerter) alertMessages() {
	alog.Info("alertMessages starting processor")

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
		logFields := map[string]interface{}{
			"traceID": eOut.Point.TraceId,
			"orgID":   eOut.Device.OrgId,
			"uniqID":  eOut.Point.UniqId,
			"devID":   eOut.Device.Id,
		}
		logger := alog.WithFields(logFields)

		// Retrieve org.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		org, err := ale.orgDAO.Read(ctx, eOut.Device.OrgId)
		cancel()
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "read"})
			logger.Errorf("alertMessages ale.orgDAO.Read: %v", err)

			continue
		}
		logger.Debugf("alertMessages org: %+v", org)

		// Retrieve alarms by rule ID. Alarms may be disabled.
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		alarms, _, err := ale.alarmDAO.List(ctx, eOut.Device.OrgId, time.Time{},
			"", 0, eOut.Rule.Id)
		cancel()
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "list"})
			logger.Errorf("alertMessages ale.alarmDAO.List: %v", err)

			continue
		}
		logger.Debugf("alertMessages alarms: %+v", alarms)

		// Validate, retrieve users, process and send alerts, and store results.
		// Unconditionally acknowledge a message after processing, as there are
		// no guarantees of alarms or users being assigned to an event.
		for _, a := range alarms {
			// Validate alarm.
			if a.Status != common.Status_ACTIVE {
				continue
			}

			// Retrieve users. Only active users with matching tags will be
			// returned.
			ctx, cancel := context.WithTimeout(context.Background(),
				5*time.Second)
			users, err := ale.userDAO.ListByTags(ctx, eOut.Device.OrgId,
				a.UserTags)
			cancel()
			if err != nil {
				metric.Incr("error", map[string]string{"func": "listbytags"})
				logger.Errorf("alertMessages ale.userDAO.ListByTags: %v", err)

				continue
			}
			if len(users) == 0 {
				continue
			}

			// Generate alert subject and body.
			subj, err := alarm.Generate(eOut.Point, eOut.Rule, eOut.Device,
				a.SubjectTemplate)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "gensubject"})
				logger.Errorf("alertMessages subject alarm.Generate: %v", err)

				continue
			}

			body, err := alarm.Generate(eOut.Point, eOut.Rule, eOut.Device,
				a.BodyTemplate)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "genbody"})
				logger.Errorf("alertMessages body alarm.Generate: %v", err)

				continue
			}
			metric.Incr("generated", nil)

			// Process alerts.
			for _, user := range users {
				// Check cache for existing repeat interval.
				key := repeatKey(eOut.Device.OrgId, eOut.Device.Id, a.Id,
					user.Id)
				ctx, cancel := context.WithTimeout(context.Background(),
					5*time.Second)
				ok, err := ale.cache.SetIfNotExistTTL(ctx, key, 1,
					time.Duration(a.RepeatInterval)*time.Minute)
				cancel()
				if err != nil {
					metric.Incr("error", map[string]string{
						"func": "setifnotexist"})
					logger.Errorf("alertMessages ale.cache.SetIfNotExistTTL: "+
						"%v", err)

					continue
				}
				if !ok {
					metric.Incr("repeat", nil)

					continue
				}

				// Send alert.
				ctx, cancel = context.WithTimeout(context.Background(),
					time.Minute)
				switch a.Type {
				case api.AlarmType_APP:
					err = ale.notify.App(ctx, user.AppKey, subj, body)
				case api.AlarmType_SMS:
					err = ale.notify.SMS(ctx, user.Phone, subj, body)
				case api.AlarmType_EMAIL:
					err = ale.notify.Email(ctx, org.DisplayName, org.Email,
						user.Email, subj, body)
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
						"type": a.Type.String()})
					logger.Debugf("alertMessages sent user, msg: %+v, %v", user,
						subj+" - "+body)
				}

				// Store alert.
				ctx, cancel = context.WithTimeout(context.Background(),
					5*time.Second)
				err = ale.alertDAO.Create(ctx, alert)
				cancel()
				if err != nil {
					metric.Incr("error", map[string]string{"func": "create"})
					logger.Errorf("alertMessages ale.alertDAO.Create: %v", err)

					continue
				}

				metric.Incr("created", nil)
				logger.Debugf("alertMessages created: %+v", alert)
			}
		}

		msg.Ack()
		metric.Incr("processed", nil)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("alertMessages processed %v messages", processCount)
		}
	}
}
