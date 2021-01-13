package validator

import (
	"context"
	"errors"
	"time"

	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/proto"
)

// validateMessages validates received device messages and builds messages
// for publishing.
func (val *Validator) validateMessages() {
	alog.Info("validateMessages starting processor")

	var processCount int
	for msg := range val.vInSub.C() {
		// Retrieve published message.
		vIn := &message.ValidatorIn{}
		err := proto.Unmarshal(msg.Payload(), vIn)
		if err != nil || vIn.Point == nil {
			msg.Ack()
			alog.Errorf("validateMessages proto.Unmarshal vIn, err: %+v, %v",
				vIn, err)
			continue
		}

		// Set up logging fields.
		logFields := map[string]interface{}{
			"traceID": vIn.Point.TraceId,
			"orgID":   vIn.OrgId,
			"uniqID":  vIn.Point.UniqId,
		}
		logger := alog.WithFields(logFields)

		// Retrieve device.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dev, err := val.devDAO.ReadByUniqID(ctx, vIn.Point.UniqId)
		cancel()
		if errors.Is(err, dao.ErrNotFound) {
			logger.Debugf("validateMessages device not found: %+v", vIn)
			msg.Ack()
			continue
		}
		if err != nil {
			logger.Errorf("validateMessages val.devDAO.ReadByUniqID: %v", err)
			msg.Requeue()
			continue
		}
		logger = logger.WithStr("devID", dev.Id)

		// Perform validation.
		switch err := vIn.Point.Validate(); {
		case err != nil:
			logger.Debugf("validateMessages vIn.Point.Validate: %v", err)
			msg.Ack()
			continue
		case vIn.OrgId != dev.OrgId:
			logger.Errorf("validateMessages incorrect org ID, expected: %v, "+
				"actual: %v", dev.OrgId, vIn.OrgId)
			msg.Ack()
			continue
		case dev.Status != common.Status_ACTIVE:
			logger.Debugf("validateMessages device disabled: %+v", vIn)
			msg.Ack()
			continue
		case vIn.Point.Token != dev.Token:
			logger.Debugf("validateMessages invalid token: %+v", vIn)
			msg.Ack()
			continue
		}

		// Build and publish ValidatorOut message.
		vOut := &message.ValidatorOut{
			Point: vIn.Point,
			OrgId: vIn.OrgId,
			DevId: dev.Id,
		}

		bVOut, err := proto.Marshal(vOut)
		if err != nil {
			logger.Errorf("validateMessages proto.Marshal: %v", err)
			msg.Ack()
			continue
		}

		if err = val.vOutQueue.Publish(val.vOutPubTopic, bVOut); err != nil {
			logger.Errorf("validateMessages val.pub.Publish: %v", err)
			msg.Requeue()
			continue
		}
		msg.Ack()
		logger.Debugf("validateMessages published: %+v", vOut)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("validateMessages processed %v messages", processCount)
		}
	}
}
