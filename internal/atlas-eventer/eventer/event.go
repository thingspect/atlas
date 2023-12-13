package eventer

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/rule"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/proto"
)

// eventMessages receives data point messages, events based on rule processing,
// and builds messages for publishing.
func (ev *Eventer) eventMessages() {
	alog.Info("eventMessages starting processor")
	ctx := context.Background()

	var processCount int
	for msg := range ev.vOutSub.C() {
		// Retrieve published message.
		metric.Incr("received", nil)
		vOut := &message.ValidatorOut{}
		err := proto.Unmarshal(msg.Payload(), vOut)
		if err != nil || vOut.GetPoint() == nil || vOut.GetDevice() == nil {
			msg.Ack()

			if !bytes.Equal([]byte{queue.Prime}, msg.Payload()) {
				metric.Incr("error", map[string]string{"func": "unmarshal"})
				alog.Errorf("eventMessages proto.Unmarshal vOut, err: %+v, %v",
					vOut, err)
			}

			continue
		}

		// Set up logging fields.
		logger := alog.
			WithField("traceID", vOut.GetPoint().GetTraceId()).
			WithField("orgID", vOut.GetDevice().GetOrgId()).
			WithField("uniqID", vOut.GetPoint().GetUniqId()).
			WithField("devID", vOut.GetDevice().GetId())

		// Retrieve rules. Only active rules with matching tags and attributes
		// will be returned.
		dCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		rules, err := ev.ruleDAO.ListByTags(dCtx, vOut.GetDevice().GetOrgId(),
			vOut.GetPoint().GetAttr(), vOut.GetDevice().GetTags())
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
			ev.evalRules(alog.NewContext(ctx, &alog.CtxLogger{Logger: logger}),
				vOut, r)
		}

		msg.Ack()
		metric.Incr("processed", nil)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("eventMessages processed %v messages", processCount)
		}
	}
}

// evalRules evaluates rules, generates events, and optionally publishes
// EventerOut messages.
func (ev *Eventer) evalRules(
	ctx context.Context, vOut *message.ValidatorOut, r *api.Rule,
) {
	logger := alog.FromContext(ctx)

	res, err := rule.Eval(vOut.GetPoint(), r.GetExpr())
	if err != nil {
		metric.Incr("error", map[string]string{"func": "eval"})
		logger.Errorf("eventMessages rule.Eval: %v", err)

		return
	}
	metric.Incr("evaluated", map[string]string{
		"result": strconv.FormatBool(res),
	})

	if res {
		event := &api.Event{
			OrgId:     vOut.GetDevice().GetOrgId(),
			UniqId:    vOut.GetDevice().GetUniqId(),
			RuleId:    r.GetId(),
			CreatedAt: vOut.GetPoint().GetTs(),
			TraceId:   vOut.GetPoint().GetTraceId(),
		}

		dCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := ev.evDAO.Create(dCtx, event)
		cancel()
		// Use a duplicate event as a tombstone to protect against failure
		// mid-loop and support fast-forward. Do not attempt to coordinate event
		// success with publish failures.
		if errors.Is(err, dao.ErrAlreadyExists) {
			metric.Incr("duplicate", nil)
			logger.Infof("eventMessages duplicate ev.evDAO.Create: %v", err)

			return
		}
		if err != nil {
			metric.Incr("error", map[string]string{"func": "create"})
			logger.Errorf("eventMessages ev.evDAO.Create: %v", err)

			return
		}

		eOut := &message.EventerOut{
			Point:  vOut.GetPoint(),
			Device: vOut.GetDevice(),
			Rule:   r,
		}
		bEOut, err := proto.Marshal(eOut)
		if err != nil {
			metric.Incr("error", map[string]string{"func": "marshal"})
			logger.Errorf("eventMessages proto.Marshal: %v", err)

			return
		}

		if err = ev.evQueue.Publish(ev.eOutPubTopic,
			bEOut); err != nil {
			metric.Incr("error", map[string]string{"func": "publish"})
			logger.Errorf("eventMessages ev.evQueue.Publish: %v", err)

			return
		}

		metric.Incr("published", nil)
		logger.Debugf("eventMessages published: %+v", eOut)
	}
}
