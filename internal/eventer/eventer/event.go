package eventer

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/rule"
	"google.golang.org/protobuf/proto"
)

// eventMessages receives data point messages, events based on rule processing,
// and builds messages for publishing.
func (ev *Eventer) eventMessages() {
	alog.Info("eventMessages starting processor")

	var processCount int
	for msg := range ev.vOutSub.C() {
		// Retrieve published message.
		metric.Incr("received", nil)
		vOut := &message.ValidatorOut{}
		err := proto.Unmarshal(msg.Payload(), vOut)
		if err != nil || vOut.Point == nil || vOut.Device == nil {
			msg.Ack()
			metric.Incr("error", map[string]string{"func": "unmarshal"})
			alog.Errorf("eventMessages proto.Unmarshal vOut, err: %+v, %v",
				vOut, err)

			continue
		}

		// Set up logging fields.
		logFields := map[string]interface{}{
			"traceID": vOut.Point.TraceId,
			"orgID":   vOut.Device.OrgId,
			"uniqID":  vOut.Point.UniqId,
			"devID":   vOut.Device.Id,
		}
		logger := alog.WithFields(logFields)

		// Retrieve rules. Only active rules with matching tags and attributes
		// will be returned.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		rules, err := ev.ruleDAO.ListByTags(ctx, vOut.Device.OrgId,
			vOut.Point.Attr, vOut.Device.Tags)
		cancel()
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "listbytags"})
			logger.Errorf("eventMessages ev.ruleDAO.ListByTags: %v", err)

			continue
		}
		logger.Debugf("eventMessages rules: %+v", rules)

		// Evaluate, event, and optionally publish EventerOut messages.
		for _, r := range rules {
			res, err := rule.Eval(vOut.Point, r.Expr)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "eval"})
				logger.Errorf("eventMessages rule.Eval: %v", err)

				continue
			}
			metric.Incr("evaluated", map[string]string{
				"result": strconv.FormatBool(res)})

			if res {
				event := &api.Event{
					OrgId:     vOut.Device.OrgId,
					UniqId:    vOut.Device.UniqId,
					RuleId:    r.Id,
					CreatedAt: vOut.Point.Ts,
					TraceId:   vOut.Point.TraceId,
				}

				ctx, cancel := context.WithTimeout(context.Background(),
					5*time.Second)
				err := ev.eventDAO.Create(ctx, event)
				cancel()
				// Use a duplicate event as a tombstone to protect against
				// failure mid-loop and support fast-forward. Do not attempt to
				// coordinate event success with publish failures.
				if errors.Is(err, dao.ErrAlreadyExists) {
					metric.Incr("duplicate", nil)
					logger.Infof("eventMessages duplicate ev.eventDAO.Create: "+
						"%v", err)

					continue
				}
				if err != nil {
					metric.Incr("error", map[string]string{"func": "create"})
					logger.Errorf("eventMessages ev.eventDAO.Create: %v", err)

					continue
				}

				eOut := &message.EventerOut{
					Point:  vOut.Point,
					Device: vOut.Device,
					Rule:   r,
				}
				bEOut, err := proto.Marshal(eOut)
				if err != nil {
					metric.Incr("error", map[string]string{"func": "marshal"})
					logger.Errorf("eventMessages proto.Marshal: %v", err)

					continue
				}

				if err = ev.evQueue.Publish(ev.eOutPubTopic,
					bEOut); err != nil {
					metric.Incr("error", map[string]string{"func": "publish"})
					logger.Errorf("eventMessages ev.evQueue.Publish: %v", err)

					continue
				}

				metric.Incr("published", nil)
				logger.Debugf("eventMessages published: %+v", eOut)
			}
		}

		msg.Ack()
		metric.Incr("processed", nil)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("eventMessages processed %v messages", processCount)
		}
	}
}
