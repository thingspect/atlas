package validator

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"google.golang.org/protobuf/proto"
)

// validateMessages validates received device messages and builds messages
// for publishing.
func (val *Validator) validateMessages() {
	alog.Info("validateMessages starting processor")

	var processCount int
	for msg := range val.vInSub.C() {
		// Retrieve published message.
		metric.Incr("received", nil)
		vIn := &message.ValidatorIn{}
		err := proto.Unmarshal(msg.Payload(), vIn)
		if err != nil || vIn.GetPoint() == nil {
			msg.Ack()

			if !bytes.Equal([]byte{queue.Prime}, msg.Payload()) {
				metric.Incr("error", map[string]string{"func": "unmarshal"})
				alog.Errorf("validateMessages proto.Unmarshal vIn, err: %+v, "+
					"%v", vIn, err)
			}

			continue
		}

		// Set up logging fields.
		logger := alog.
			WithField("traceID", vIn.GetPoint().GetTraceId()).
			WithField("orgID", vIn.GetOrgId()).
			WithField("uniqID", vIn.GetPoint().GetUniqId())

		// Retrieve device.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dev, err := val.devDAO.ReadByUniqID(ctx, vIn.GetPoint().GetUniqId())
		cancel()
		if errors.Is(err, dao.ErrNotFound) {
			msg.Ack()
			metric.Incr("notfound", nil)
			logger.Debugf("validateMessages device not found: %+v", vIn)

			continue
		}
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "readbyuniqid"})
			logger.Errorf("validateMessages val.devDAO.ReadByUniqID: %v", err)

			continue
		}
		logger = logger.WithField("devID", dev.GetId())

		// Perform validation.
		switch err := vIn.GetPoint().Validate(); {
		case err != nil:
			msg.Ack()
			metric.Incr("invalid", map[string]string{"func": "validate"})
			logger.Debugf("validateMessages vIn.Point.Validate: %v", err)

			continue
		case !vIn.GetSkipToken() && vIn.GetOrgId() != dev.GetOrgId():
			msg.Ack()
			metric.Incr("invalid", map[string]string{"func": "orgid"})
			logger.Errorf("validateMessages incorrect org ID, expected: %v, "+
				"actual: %v", dev.GetOrgId(), vIn.GetOrgId())

			continue
		case dev.GetStatus() != api.Status_ACTIVE:
			msg.Ack()
			metric.Incr("invalid", map[string]string{"func": "disabled"})
			logger.Debugf("validateMessages device disabled: %+v", vIn)

			continue
		case !vIn.GetSkipToken() && vIn.GetPoint().GetToken() != dev.GetToken():
			msg.Ack()
			metric.Incr("invalid", map[string]string{"func": "token"})
			logger.Debugf("validateMessages invalid token: %+v", vIn)

			continue
		}
		metric.Incr("processed", nil)

		// Build and publish ValidatorOut message.
		vOut := &message.ValidatorOut{
			Point:  vIn.GetPoint(),
			Device: dev,
		}

		bVOut, err := proto.Marshal(vOut)
		if err != nil {
			msg.Ack()
			metric.Incr("error", map[string]string{"func": "marshal"})
			logger.Errorf("validateMessages proto.Marshal: %v", err)

			continue
		}

		if err = val.valQueue.Publish(val.vOutPubTopic, bVOut); err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "publish"})
			logger.Errorf("validateMessages val.pub.Publish: %v", err)

			continue
		}

		msg.Ack()
		metric.Incr("published", nil)
		logger.Debugf("validateMessages published: %+v", vOut)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("validateMessages processed %v messages", processCount)
		}
	}
}
