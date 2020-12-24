package accumulator

import (
	"context"
	"errors"
	"time"

	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/proto"
)

// accumulateMessages accumulates data point messages and stores them.
func (acc *Accumulator) accumulateMessages() {
	alog.Info("accumulateMessages starting processor")

	var processCount int
	for msg := range acc.vOutSub.C() {
		// Retrieve published message.
		vOut := &message.ValidatorOut{}
		err := proto.Unmarshal(msg.Payload(), vOut)
		if err != nil || vOut.Point == nil {
			msg.Ack()
			alog.Errorf("validateMessages proto.Unmarshal vOut, err: %+v, %v",
				vOut, err)
			continue
		}

		// Set up logging fields.
		logFields := map[string]interface{}{
			"traceID": vOut.Point.TraceId,
			"orgID":   vOut.OrgId,
			"uniqID":  vOut.Point.UniqId,
			"devID":   vOut.DevId,
		}
		logEntry := alog.WithFields(logFields)

		// Create data point.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = acc.dpDAO.Create(ctx, vOut.Point, vOut.OrgId)
		cancel()
		if errors.Is(err, dao.ErrAlreadyExists) ||
			errors.Is(err, dao.ErrInvalidFormat) {
			logEntry.Errorf("accumulateMessages discard acc.dpDAO.Create: %v",
				err)
			msg.Ack()
			continue
		}
		if err != nil {
			logEntry.Errorf("accumulateMessages requeue acc.dpDAO.Create: %v",
				err)
			msg.Requeue()
			continue
		}

		msg.Ack()
		logEntry.Debugf("accumulateMessages created: %+v", vOut)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("accumulateMessages processed %v messages", processCount)
		}
	}
}
