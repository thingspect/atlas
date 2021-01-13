package ingestor

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/api/go/mqtt"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// parseMessages parses MQTT messages and builds messages for publishing.
func (ing *Ingestor) parseMessages() {
	alog.Info("parseMessages starting processor")

	var processCount int
	for msg := range ing.mqttSub.C() {
		// Set up logging fields.
		traceID := uuid.New().String()
		logger := alog.WithStr("traceID", traceID)

		// Parse and validate topic in format: 'v1/:orgID[/:uniqID][/json]'.
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		if len(topicParts) < 2 || len(topicParts) > 4 || topicParts[0] != "v1" {
			msg.Ack()
			logger.Errorf("parseMessages malformed topic: %v", topic)
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
			msg.Ack()
			logger.Errorf("parseMessages proto.Unmarshal: %v", err)
			continue
		}
		logger.Debugf("parseMessages payl: %+v", payl)

		// Build and publish ValidatorIn messages.
		var successCount int
		for _, point := range payl.Points {
			vIn := dataPointToVIn(traceID, payl.Token, topicParts, point)

			bVIn, err := proto.Marshal(vIn)
			if err != nil {
				logger.Errorf("parseMessages proto.Marshal: %v", err)
				continue
			}

			if err = ing.parserQueue.Publish(ing.parserPubTopic,
				bVIn); err != nil {
				logger.Errorf("parseMessages ing.parserPub.Publish: %v", err)
				continue
			}

			successCount++
			logger.Debugf("parseMessages published: %+v", vIn)
		}

		// Do not ack on errors, as publish may retry successfully.
		// Deduplication will take place downstream.
		if successCount == len(payl.Points) {
			msg.Ack()
		}

		processCount++
		if processCount%100 == 0 {
			alog.Infof("parseMessages processed %v messages", processCount)
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
