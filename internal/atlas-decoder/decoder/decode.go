package decoder

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/decode/registry"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/proto/go/message"
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

			if !bytes.Equal([]byte{queue.Prime}, msg.Payload()) {
				metric.Incr("error", map[string]string{"func": "unmarshal"})
				alog.Errorf("decodeMessages proto.Unmarshal dIn, err: %+v, %v",
					dIn, err)
			}

			continue
		}

		// Trace IDs have been authenticated and are safe to copy.
		var traceID uuid.UUID
		copy(traceID[:], dIn.GetTraceId())

		// Set up logging fields.
		logger := alog.
			WithField("traceID", traceID.String()).
			WithField("uniqID", dIn.GetUniqId())

		// Retrieve device.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		dev, err := dec.devDAO.ReadByUniqID(ctx, dIn.GetUniqId())
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
		logger = logger.WithField("orgID", dev.GetOrgId())
		logger = logger.WithField("devID", dev.GetId())

		// Decode data payload. Continue execution in the presence of errors, as
		// valid points may be returned.
		points, err := dec.reg.Decode(dev.GetDecoder(), dIn.GetData())
		if err != nil {
			metric.Incr("error", map[string]string{"func": "decode"})
			logger.Errorf("decodeMessages dec.registry.Decode: %v", err)
		}
		metric.Incr("processed", nil)
		logger.Debugf("decodeMessages points: %+v", points)

		// Build and publish ValidatorIn messages.
		var successCount int
		for _, point := range points {
			vIn := registry.PointToVIn(traceID.String(), dIn.GetUniqId(), point,
				dIn.GetTs())

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("decodeMessages proto.Marshal: %v", err)

				continue
			}

			if err = dec.decQueue.Publish(dec.vInPubTopic,
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

		// Do not ack on points loop errors, as publish may retry successfully.
		// Deduplication will take place downstream.
		if successCount == len(points) {
			msg.Ack()
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("decodeMessages processed %v messages", processCount)
		}
	}
}
