package ingestor

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/api/go/mqtt"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// decodeMessages decodes MQTT messages and builds messages for publishing.
func (ing *Ingestor) decodeMessages() {
	alog.Info("decodeMessages starting processor")

	var processCount int
	for msg := range ing.mqttSub.C() {
		msg.Ack()
		metric.Incr("received", nil)

		// Set up logging fields.
		traceID := uuid.NewString()
		logger := alog.WithStr("traceID", traceID)

		// Parse and validate topic in format: 'v1/:orgID[/:uniqID][/json]'.
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		if len(topicParts) < 2 || len(topicParts) > 4 || topicParts[0] != "v1" {
			metric.Incr("error", map[string]string{"func": "topic"})
			logger.Errorf("decodeMessages malformed topic: %v", topic)
			continue
		}
		logger = logger.WithStr("orgID", topicParts[1])

		// Unmarshal payload based on topic and format.
		payl := &mqtt.Payload{}
		var err error

		if topicParts[len(topicParts)-1] == "json" {
			logger = logger.WithStr("paylType", "json")
			topicParts = topicParts[:len(topicParts)-1]
			err = protojson.Unmarshal(msg.Payload(), payl)
		} else {
			logger = logger.WithStr("paylType", "proto")
			err = proto.Unmarshal(msg.Payload(), payl)
		}
		if err != nil {
			metric.Incr("error", map[string]string{"func": "unmarshal"})
			logger.Errorf("decodeMessages proto.Unmarshal: %v", err)
			continue
		}
		metric.Incr("processed", nil)
		logger.Debugf("decodeMessages payl: %+v", payl)

		// Build and publish ValidatorIn messages.
		var successCount int
		for _, point := range payl.Points {
			vIn := dataPointToVIn(traceID, payl.Token, topicParts, point)

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				metric.Incr("error", map[string]string{"func": "marshal"})
				logger.Errorf("decodeMessages proto.Marshal: %v", err)
				continue
			}

			if err = ing.decoderQueue.Publish(ing.decoderPubTopic,
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

		processCount++
		if processCount%100 == 0 {
			alog.Infof("decodeMessages processed %v messages", processCount)
		}
	}
}

// dataPointToVIn maps a DataPoint to ValidatorIn. The DataPoint is embedded in
// ValidatorIn to avoid copying or use of Clone/reflection. Tests should take
// this into account.
func dataPointToVIn(traceID, paylToken string, topicParts []string,
	point *common.DataPoint) *message.ValidatorIn {
	vIn := &message.ValidatorIn{
		Point: point,
		OrgId: topicParts[1],
	}

	// Override trace ID.
	vIn.Point.TraceId = traceID

	// Override UniqID with topic-based ID, if present.
	if len(topicParts) == 3 {
		vIn.Point.UniqId = topicParts[2]
	}

	// Default to current timestamp if not provided.
	if vIn.Point.Ts == nil {
		vIn.Point.Ts = timestamppb.Now()
	}

	// Override Token with payload-based Token, if present.
	if paylToken != "" {
		vIn.Point.Token = paylToken
	}
	return vIn
}
