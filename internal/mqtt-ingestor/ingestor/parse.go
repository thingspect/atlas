package ingestor

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/mqtt"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// parseMessages parses MQTT messages and builds messages for publishing.
func (ing *Ingestor) parseMessages() {
	alog.Info("parseMessages starting processor")

	var processCount int
	for msg := range ing.mqttSub.C() {
		// Set up logging fields.
		traceID := uuid.New().String()
		logEntry := alog.WithStr("traceID", traceID)

		// Parse and validate topic in format: 'v1/:orgID[/:uniqID][/json]'.
		topic := msg.Topic()
		topicParts := strings.Split(topic, "/")
		if len(topicParts) < 2 || len(topicParts) > 4 || topicParts[0] != "v1" {
			msg.Ack()
			logEntry.Errorf("parseMessages malformed topic: %v", topic)
			continue
		}
		logEntry = logEntry.WithStr("orgID", topicParts[1])

		// Unmarshal payload based on topic and format.
		payl := &mqtt.Payload{}
		var err error

		if topicParts[len(topicParts)-1] == "json" {
			logEntry = logEntry.WithStr("paylType", "json")
			topicParts = topicParts[:len(topicParts)-1]
			err = protojson.Unmarshal(msg.Payload(), payl)
		} else {
			logEntry = logEntry.WithStr("paylType", "proto")
			err = proto.Unmarshal(msg.Payload(), payl)
		}
		if err != nil {
			msg.Ack()
			logEntry.Errorf("parseMessages proto.Unmarshal: %v", err)
			continue
		}
		logEntry.Debugf("parseMessages payl: %#v", payl)

		// Build and publish MQTTIngestorOut message.
		var successCount int
		for _, point := range payl.Points {
			mIOut := dataPointToMIOut(traceID, payl.Token, topicParts, point)

			// Marshal MQTTIngestorOut.
			bMIOut, err := proto.Marshal(mIOut)
			if err != nil {
				logEntry.Errorf("parseMessages proto.Marshal: %v", err)
				continue
			}

			// Publish message.
			if err = ing.parserPub.Publish(pubTopic, bMIOut); err != nil {
				logEntry.Errorf("parseMessages ing.parserPub.Publish: %v", err)
				continue
			}
			successCount++
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

// dataPointToMIOut maps a DataPoint to MQTTIngestorOut.
func dataPointToMIOut(traceID, token string, topicParts []string,
	point *mqtt.DataPoint) *message.MQTTIngestorOut {
	mIOut := &message.MQTTIngestorOut{}
	mIOut.UniqId = point.UniqId
	mIOut.Attr = point.Attr

	// If none of the types match, it is a map or absent.
	switch point.ValOneof.(type) {
	case *mqtt.DataPoint_IntVal:
		mIOut.ValOneof = &message.MQTTIngestorOut_IntVal{
			IntVal: point.GetIntVal(),
		}
	case *mqtt.DataPoint_Fl64Val:
		mIOut.ValOneof = &message.MQTTIngestorOut_Fl64Val{
			Fl64Val: point.GetFl64Val(),
		}
	case *mqtt.DataPoint_StrVal:
		mIOut.ValOneof = &message.MQTTIngestorOut_StrVal{
			StrVal: point.GetStrVal(),
		}
	case *mqtt.DataPoint_BoolVal:
		mIOut.ValOneof = &message.MQTTIngestorOut_BoolVal{
			BoolVal: point.GetBoolVal(),
		}
	}

	mIOut.MapVal = point.MapVal
	mIOut.Ts = point.Ts
	mIOut.Token = token
	mIOut.OrgId = topicParts[1]
	mIOut.TraceId = traceID

	// Override UniqID with topic-based ID, if present.
	if len(topicParts) == 3 {
		mIOut.UniqId = topicParts[2]
	}
	return mIOut
}
