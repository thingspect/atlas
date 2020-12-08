package validator

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
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
		logEntry := alog.WithFields(logFields)

		// Retrieve device and begin validation.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dev, err := val.devDAO.ReadByUniqID(ctx, vIn.Point.UniqId)
		cancel()
		if errors.Is(err, sql.ErrNoRows) {
			logEntry.Debugf("validateMessages device not found: %+v", vIn)
			msg.Ack()
			continue
		}
		if err != nil {
			logEntry.Errorf("validateMessages val.devDAO.ReadByUniqID: %v", err)
			msg.Requeue()
			continue
		}
		logEntry = logEntry.WithStr("devID", dev.ID)

		switch {
		case vIn.OrgId != dev.OrgID:
			logEntry.Debugf("validateMessages invalid org ID: %+v", vIn)
			msg.Ack()
			continue
		case dev.Disabled:
			logEntry.Debugf("validateMessages device disabled: %+v", vIn)
			msg.Ack()
			continue
		case vIn.Point.Token != dev.Token:
			logEntry.Debugf("validateMessages invalid token: %+v", vIn)
			msg.Ack()
			continue
		}

		// Build and publish ValidatorOut message.
		vOut := &message.ValidatorOut{
			Point: vIn.Point,
			OrgId: vIn.OrgId,
			DevId: dev.ID,
		}

		bVOut, err := proto.Marshal(vOut)
		if err != nil {
			logEntry.Errorf("validateMessages proto.Marshal: %v", err)
			msg.Ack()
			continue
		}

		if err = val.vOutQueue.Publish(val.vOutPubTopic, bVOut); err != nil {
			logEntry.Errorf("validateMessages val.pub.Publish: %v", err)
			msg.Requeue()
			continue
		}
		msg.Ack()
		logEntry.Debugf("validateMessages published: %#v", vOut)

		processCount++
		if processCount%100 == 0 {
			alog.Infof("validateMessages processed %v messages", processCount)
		}
	}
}
