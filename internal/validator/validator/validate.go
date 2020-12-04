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
		if err != nil {
			msg.Ack()
			alog.Errorf("validateMessages proto.Unmarshal: %v", err)
			continue
		}

		// Set up logging fields.
		logFields := map[string]interface{}{
			"traceID": vIn.TraceId,
			"orgID":   vIn.OrgId,
			"uniqID":  vIn.UniqId,
		}
		logEntry := alog.WithFields(logFields)

		// Retrieve device and begin validation.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dev, err := val.devDAO.ReadByUniqID(ctx, vIn.UniqId)
		cancel()
		if errors.Is(err, sql.ErrNoRows) {
			logEntry.Debugf("validateMessages device not found: %#v", vIn)
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
			logEntry.Debugf("validateMessages invalid org ID: %#v", vIn)
			msg.Ack()
			continue
		case dev.Disabled:
			logEntry.Debugf("validateMessages device disabled: %#v", vIn)
			msg.Ack()
			continue
		case vIn.Token != dev.Token:
			logEntry.Debugf("validateMessages invalid token: %#v", vIn)
			msg.Ack()
			continue
		}

		// Build and publish ValidatorOut message.
		vOut := vInToVOut(vIn, dev.ID)

		// Marshal ValidatorOut.
		bVOut, err := proto.Marshal(vOut)
		if err != nil {
			logEntry.Errorf("validateMessages proto.Marshal: %v", err)
			msg.Ack()
			continue
		}

		// Publish message.
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

// vInToVOut maps a ValidatorIn to ValidatorOut.
func vInToVOut(vIn *message.ValidatorIn, devID string) *message.ValidatorOut {
	vOut := &message.ValidatorOut{}
	vOut.UniqId = vIn.UniqId
	vOut.Attr = vIn.Attr

	// If none of the types match, it is a map or absent.
	switch vIn.ValOneof.(type) {
	case *message.ValidatorIn_IntVal:
		vOut.ValOneof = &message.ValidatorOut_IntVal{
			IntVal: vIn.GetIntVal(),
		}
	case *message.ValidatorIn_Fl64Val:
		vOut.ValOneof = &message.ValidatorOut_Fl64Val{
			Fl64Val: vIn.GetFl64Val(),
		}
	case *message.ValidatorIn_StrVal:
		vOut.ValOneof = &message.ValidatorOut_StrVal{
			StrVal: vIn.GetStrVal(),
		}
	case *message.ValidatorIn_BoolVal:
		vOut.ValOneof = &message.ValidatorOut_BoolVal{
			BoolVal: vIn.GetBoolVal(),
		}
	case *message.ValidatorIn_BytesVal:
		vOut.ValOneof = &message.ValidatorOut_BytesVal{
			BytesVal: vIn.GetBytesVal(),
		}
	}

	vOut.MapVal = vIn.MapVal
	vOut.Ts = vIn.Ts
	vOut.DevId = devID
	vOut.OrgId = vIn.OrgId
	vOut.TraceId = vIn.TraceId
	return vOut
}
