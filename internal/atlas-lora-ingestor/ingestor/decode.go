package ingestor

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/decode/chirpstack/device"
	"github.com/thingspect/atlas/pkg/decode/chirpstack/gateway"
	"github.com/thingspect/atlas/pkg/decode/registry"
	"github.com/thingspect/atlas/pkg/metric"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// decodeGateways decodes gateway messages and builds messages for publishing.
func (ing *Ingestor) decodeGateways() {
	alog.Info("decodeGateways starting processor")

	var processCount int
	for msg := range ing.mqttGWSub.C() {
		msg.Ack()
		metric.Incr("received", map[string]string{"type": "gateway"})

		// Set up logging fields.
		traceID := uuid.NewString()
		logFields := map[string]interface{}{
			"type":    "gateway",
			"traceID": traceID,
		}
		logger := alog.WithFields(logFields)

		// Parse and validate topic in formats: 'lora/+/gateway/+/event/+' and
		// 'lora/+/gateway/+/state/+'. A state type is considered an event.
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		if len(topicParts) != 6 ||
			topicParts[0] != "lora" ||
			topicParts[2] != "gateway" ||
			(topicParts[4] != "event" && topicParts[4] != "state") {
			metric.Incr("error", map[string]string{"func": "topic"})
			logger.Errorf("decodeGateways unknown topic: %v", topic)

			continue
		}
		logger = logger.WithStr("uniqID", topicParts[3])
		logger = logger.WithStr("event", topicParts[5])

		// Decode payload. Continue execution in the presence of errors, as
		// valid points may be returned.
		points, err := gateway.Parse(topicParts[5], msg.Payload())
		if err != nil {
			metric.Incr("error", map[string]string{"func": "decode"})
			logger.Errorf("decodeGateways gateway.Gateway: %v", err)
		}
		metric.Incr("processed", map[string]string{"type": "gateway"})
		logger.Debugf("decodeGateways points: %+v", points)

		// Build and publish ValidatorIn messages.
		for _, point := range points {
			vIn := registry.PointToVIn(traceID, topicParts[3], point,
				timestamppb.Now())

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("decodeGateways proto.Marshal: %v", err)

				continue
			}

			if err = ing.ingQueue.Publish(ing.vInGWPubTopic,
				bVIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("decodeGateways ing.decoderQueue.Publish: %v",
					err)

				continue
			}

			metric.Incr("published", map[string]string{"type": "gateway"})
			logger.Debugf("decodeGateways published: %+v", vIn)
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("decodeGateways processed %v messages", processCount)
		}
	}
}

// decodeDevices decodes device messages and builds messages for publishing.
func (ing *Ingestor) decodeDevices() {
	alog.Info("decodeDevices starting processor")

	var processCount int
	for msg := range ing.mqttDevSub.C() {
		msg.Ack()
		metric.Incr("received", map[string]string{"type": "device"})

		// Set up logging fields.
		traceID := uuid.NewString()
		logFields := map[string]interface{}{
			"type":    "device",
			"traceID": traceID,
		}
		logger := alog.WithFields(logFields)

		// Parse and validate topic in format:
		// 'lora/application/+/device/+/event/+'.
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		if len(topicParts) != 7 || topicParts[0] != "lora" ||
			topicParts[1] != "application" || topicParts[3] != "device" ||
			topicParts[5] != "event" {
			metric.Incr("error", map[string]string{"func": "topic"})
			logger.Errorf("decodeDevices unknown topic: %v", topic)

			continue
		}
		logger = logger.WithStr("uniqID", topicParts[4])
		logger = logger.WithStr("event", topicParts[6])

		// Decode payload. Continue execution in the presence of errors, as
		// valid points may be included.
		points, ts, data, err := device.Parse(topicParts[6], msg.Payload())
		if err != nil {
			metric.Incr("error", map[string]string{"func": "decode"})
			logger.Errorf("decodeDevices device.Device: %v", err)
		}
		metric.Incr("processed", map[string]string{"type": "device"})
		logger.Debugf("decodeDevices points: %+v", points)

		// Build and publish ValidatorIn messages.
		for _, point := range points {
			vIn := registry.PointToVIn(traceID, topicParts[4], point, ts)

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("decodeDevices point proto.Marshal: %v", err)

				continue
			}

			if err = ing.ingQueue.Publish(ing.vInDevPubTopic,
				bVIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("decodeDevices point ing.decoderQueue.Publish: "+
					"%v", err)

				continue
			}

			metric.Incr("published", map[string]string{"type": "device"})
			logger.Debugf("decodeDevices point published: %+v", vIn)
		}

		// Build and publish DecoderIn messages, if present.
		if len(data) > 0 {
			pIn := &message.DecoderIn{
				UniqId:  topicParts[4],
				Data:    data,
				Ts:      ts,
				TraceId: traceID,
			}

			bPIn, err := proto.Marshal(pIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("decodeDevices data proto.Marshal: %v", err)

				continue
			}

			if err = ing.ingQueue.Publish(ing.dInPubTopic,
				bPIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("decodeDevices data ing.decoderQueue.Publish: %v",
					err)

				continue
			}

			metric.Incr("published", map[string]string{"type": "data"})
			logger.Debugf("decodeDevices data published: %+v", pIn)
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("decodeDevices processed %v messages", processCount)
		}
	}
}
