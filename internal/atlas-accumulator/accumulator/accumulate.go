package accumulator

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"google.golang.org/protobuf/proto"
)

// accumulateMessages accumulates data point messages and stores them.
func (acc *Accumulator) accumulateMessages() {
	alog.Info("accumulateMessages starting processor")

	var processCount int
	for msg := range acc.vOutSub.C() {
		// Retrieve published message.
		metric.Incr("received", nil)
		vOut := &message.ValidatorOut{}
		err := proto.Unmarshal(msg.Payload(), vOut)
		if err != nil || vOut.Point == nil || vOut.Device == nil {
			msg.Ack()

			if !bytes.Equal([]byte{queue.Prime}, msg.Payload()) {
				metric.Incr("error", map[string]string{"func": "unmarshal"})
				alog.Errorf("accumulateMessages proto.Unmarshal vOut, err: "+
					"%+v, %v", vOut, err)
			}

			continue
		}

		// Set up logging fields.
		logger := alog.
			WithField("traceID", vOut.Point.TraceId).
			WithField("orgID", vOut.Device.OrgId).
			WithField("uniqID", vOut.Point.UniqId).
			WithField("devID", vOut.Device.Id)

		// Create data point.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = acc.dpDAO.Create(ctx, vOut.Point, vOut.Device.OrgId)
		cancel()
		if errors.Is(err, dao.ErrAlreadyExists) {
			msg.Ack()
			metric.Incr("duplicate", nil)
			logger.Infof("accumulateMessages discard acc.dpDAO.Create: %v", err)

			continue
		}
		if errors.Is(err, dao.ErrInvalidFormat) {
			msg.Ack()
			metric.Incr("error", map[string]string{"func": "create"})
			logger.Errorf("accumulateMessages invalid acc.dpDAO.Create: %v",
				err)

			continue
		}
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "create"})
			logger.Errorf("accumulateMessages requeue acc.dpDAO.Create: %v",
				err)

			continue
		}

		msg.Ack()
		metric.Incr("processed", nil)
		logger.Debugf("accumulateMessages processed: %+v", vOut)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("accumulateMessages processed %v messages", processCount)
		}
	}
}
