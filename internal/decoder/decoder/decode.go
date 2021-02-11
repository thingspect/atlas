package decoder

import (
	"context"
	"errors"
	"time"

	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/decode/registry"
	"github.com/thingspect/atlas/pkg/metric"
	"google.golang.org/protobuf/proto"
)

// decodeMessages decodes received device data payload messages and builds
// messages for publishing.
func (dec *Decoder) decodeMessages() {
	alog.Info("decodeMessages starting processor")

	var processCount int
	for msg := range dec.dInSub.C() {
		// Retrieve published message.
		metric.Incr("received", nil)
		dIn := &message.DecoderIn{}
		err := proto.Unmarshal(msg.Payload(), dIn)
		if err != nil {
			msg.Ack()
			metric.Incr("error", map[string]string{"func": "unmarshal"})
			alog.Errorf("decodeMessages proto.Unmarshal dIn, err: %+v, %v",
				dIn, err)
			continue
		}

		// Set up logging fields.
		logFields := map[string]interface{}{
			"traceID": dIn.TraceId,
			"uniqID":  dIn.UniqId,
		}
		logger := alog.WithFields(logFields)

		// Retrieve device.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dev, err := dec.devDAO.ReadByUniqID(ctx, dIn.UniqId)
		cancel()
		if errors.Is(err, dao.ErrNotFound) {
			msg.Ack()
			metric.Incr("notfound", nil)
			logger.Debugf("decodeMessages device not found: %+v", dIn)
			continue
		}
		if err != nil {
			msg.Requeue()
			metric.Incr("error", map[string]string{"func": "readbyuniqid"})
			logger.Errorf("decodeMessages dec.devDAO.ReadByUniqID: %v", err)
			continue
		}
		logger = logger.WithStr("orgID", dev.OrgId)
		logger = logger.WithStr("devID", dev.Id)

		// Decode data payload. Continue execution in the presence of errors, as
		// valid points may be returned.
		points, err := dec.registry.Decode(dev.Decoder, dIn.Data)
		if err != nil {
			metric.Incr("error", map[string]string{"func": "decode"})
			logger.Errorf("decodeMessages dec.registry.Decode: %v", err)
		}
		logger.Debugf("decodeMessages points: %+v", points)

		// Build and publish ValidatorIn messages.
		var successCount int
		for _, point := range points {
			vIn := registry.PointToVIn(dIn.TraceId, dIn.UniqId, point, dIn.Ts)

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("decodeMessages proto.Marshal: %v", err)
				continue
			}

			if err = dec.decoderQueue.Publish(dec.decoderPubTopic,
				bVIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("decodeMessages ing.decoderQueue.Publish: %v",
					err)
				continue
			}

			successCount++
			metric.Incr("published", nil)
			logger.Debugf("decodeMessages published: %+v", vIn)
		}

		// Do not ack on errors, as publish may retry successfully.
		// Deduplication will take place downstream.
		if successCount == len(points) {
			msg.Ack()
			metric.Incr("processed", nil)
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("decodeMessages processed %v messages", processCount)
		}
	}
}
