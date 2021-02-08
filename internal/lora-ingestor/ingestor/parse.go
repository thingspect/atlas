package ingestor

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/parse"
	"github.com/thingspect/atlas/pkg/parse/chirpstack/device"
	"github.com/thingspect/atlas/pkg/parse/chirpstack/gateway"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// parseGateways parses gateway messages and builds messages for publishing.
func (ing *Ingestor) parseGateways() {
	alog.Info("parseGateways starting processor")

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

		// Parse and validate topic in format: 'lora/gateway/+/event/+'.
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		if len(topicParts) != 5 || topicParts[0] != "lora" ||
			topicParts[1] != "gateway" || topicParts[3] != "event" {
			metric.Incr("error", map[string]string{"func": "topic"})
			logger.Errorf("parseGateways malformed topic: %v", topic)
			continue
		}
		logger = logger.WithStr("uniqID", topicParts[2])
		logger = logger.WithStr("event", topicParts[4])

		// Parse payload. Continue execution in the presence of errors, as valid
		// points may be included.
		points, err := gateway.Gateway(topicParts[4], msg.Payload())
		if err != nil {
			metric.Incr("error", map[string]string{"func": "parse"})
			logger.Errorf("parseGateways gateway.Gateway: %v", err)
		}
		metric.Incr("processed", map[string]string{"type": "gateway"})
		logger.Debugf("parseGateways points: %+v", points)

		// Build and publish ValidatorIn messages.
		for _, point := range points {
			vIn := pointToVIn(traceID, topicParts[2], point, timestamppb.Now())

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("parseGateways proto.Marshal: %v", err)
				continue
			}

			if err = ing.parserQueue.Publish(ing.parserPubGWTopic,
				bVIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("parseGateways ing.parserPub.Publish: %v", err)
				continue
			}

			metric.Incr("published", map[string]string{"type": "gateway"})
			logger.Debugf("parseGateways published: %+v", vIn)
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("parseGateways processed %v messages", processCount)
		}
	}
}

// parseDevices parses device messages and builds messages for publishing.
func (ing *Ingestor) parseDevices() {
	alog.Info("parseDevices starting processor")

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
			logger.Errorf("parseDevices malformed topic: %v", topic)
			continue
		}
		logger = logger.WithStr("uniqID", topicParts[4])
		logger = logger.WithStr("event", topicParts[6])

		// Parse payload. Continue execution in the presence of errors, as valid
		// points may be included.
		points, ts, data, err := device.Device(topicParts[6], msg.Payload())
		if err != nil {
			metric.Incr("error", map[string]string{"func": "parse"})
			logger.Errorf("parseDevices device.Device: %v", err)
		}
		metric.Incr("processed", map[string]string{"type": "device"})
		logger.Debugf("parseDevices points: %+v", points)

		// Build and publish ValidatorIn messages.
		for _, point := range points {
			vIn := pointToVIn(traceID, topicParts[4], point, ts)

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("parseDevices point proto.Marshal: %v", err)
				continue
			}

			if err = ing.parserQueue.Publish(ing.parserPubDevTopic,
				bVIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("parseDevices point ing.parserPub.Publish: %v",
					err)
				continue
			}

			metric.Incr("published", map[string]string{"type": "device"})
			logger.Debugf("parseDevices point published: %+v", vIn)
		}

		// Build and publish ParserIn messages, if present.
		if len(data) > 0 {
			pIn := &message.ParserIn{
				UniqId:  topicParts[4],
				Data:    data,
				Ts:      ts,
				TraceId: traceID,
			}

			bPIn, err := proto.Marshal(pIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("parseDevices data proto.Marshal: %v", err)
				continue
			}

			if err = ing.parserQueue.Publish(ing.parserPubDataTopic,
				bPIn); err != nil {
				metric.Incr("error", map[string]string{"func": "publish"})
				logger.Errorf("parseDevices data ing.parserPub.Publish: %v",
					err)
				continue
			}

			metric.Incr("published", map[string]string{"type": "data"})
			logger.Debugf("parseDevices data published: %+v", pIn)
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("parseDevices processed %v messages", processCount)
		}
	}
}

// pointToVIn maps a Point to ValidatorIn.
func pointToVIn(traceID, uniqID string, point *parse.Point,
	ts *timestamppb.Timestamp) *message.ValidatorIn {
	vIn := &message.ValidatorIn{
		Point: &common.DataPoint{
			UniqId:  uniqID,
			Attr:    point.Attr,
			Ts:      ts,
			TraceId: traceID,
		},
		SkipToken: true,
	}

	switch v := point.Value.(type) {
	case int32:
		vIn.Point.ValOneof = &common.DataPoint_IntVal{IntVal: v}
	case int:
		vIn.Point.ValOneof = &common.DataPoint_IntVal{IntVal: int32(v)}
		alog.Errorf("pointToVIn casting from int: %v, %v,", point.Attr, v)
	case int64:
		vIn.Point.ValOneof = &common.DataPoint_IntVal{IntVal: int32(v)}
		alog.Errorf("pointToVIn casting from int64: %v, %v,", point.Attr, v)
	case float64:
		vIn.Point.ValOneof = &common.DataPoint_Fl64Val{Fl64Val: v}
	case float32:
		vIn.Point.ValOneof = &common.DataPoint_Fl64Val{Fl64Val: float64(v)}
		alog.Errorf("pointToVIn casting from float32: %v, %v,", point.Attr, v)
	case string:
		vIn.Point.ValOneof = &common.DataPoint_StrVal{StrVal: v}
	case bool:
		vIn.Point.ValOneof = &common.DataPoint_BoolVal{BoolVal: v}
	case []byte:
		vIn.Point.ValOneof = &common.DataPoint_BytesVal{BytesVal: v}
	default:
		alog.Errorf("pointToVIn unknown type: %v, %T, %v,", point.Attr,
			point.Value, point.Value)
	}

	return vIn
}
