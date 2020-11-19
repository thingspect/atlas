// +build !integration

package ingestor

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/mqtt"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseMessages(t *testing.T) {
	t.Parallel()

	orgID := uuid.New().String()
	token := uuid.New().String()
	uniqIDTopic := random.String(16)
	uniqIDPoint := random.String(16)
	now := timestamppb.Now()

	tests := []struct {
		inpTopic string
		inpPoint *mqtt.DataPoint
		res      *message.MQTTIngestorOut
	}{
		{fmt.Sprintf("v1/%s/%s", orgID, uniqIDTopic),
			&mqtt.DataPoint{Attr: "motion",
				ValOneof: &mqtt.DataPoint_IntVal{IntVal: 123}},
			&message.MQTTIngestorOut{UniqId: uniqIDTopic, Attr: "motion",
				ValOneof: &message.MQTTIngestorOut_IntVal{IntVal: 123},
				Token:    token, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/json", orgID), &mqtt.DataPoint{UniqId: uniqIDPoint,
			Attr: "temp", ValOneof: &mqtt.DataPoint_Fl64Val{Fl64Val: 20.3}},
			&message.MQTTIngestorOut{UniqId: uniqIDPoint, Attr: "temp",
				ValOneof: &message.MQTTIngestorOut_Fl64Val{Fl64Val: 20.3},
				Token:    token, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/%s/json", orgID, uniqIDTopic), &mqtt.DataPoint{
			Attr: "power", ValOneof: &mqtt.DataPoint_StrVal{StrVal: "batt"},
			Ts: now}, &message.MQTTIngestorOut{UniqId: uniqIDTopic,
			Attr: "power", ValOneof: &message.MQTTIngestorOut_StrVal{
				StrVal: "batt"}, Ts: now, Token: token, OrgId: orgID}},
		{fmt.Sprintf("v1/%s", orgID), &mqtt.DataPoint{UniqId: uniqIDPoint,
			Attr: "leak", ValOneof: &mqtt.DataPoint_BoolVal{BoolVal: true}},
			&message.MQTTIngestorOut{UniqId: uniqIDPoint, Attr: "leak",
				ValOneof: &message.MQTTIngestorOut_BoolVal{BoolVal: true},
				Token:    token, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/%s", orgID, uniqIDTopic), &mqtt.DataPoint{
			Attr: "metadata", MapVal: map[string]string{"aaa": "bbb"}, Ts: now},
			&message.MQTTIngestorOut{UniqId: uniqIDTopic, Attr: "metadata",
				MapVal: map[string]string{"aaa": "bbb"}, Ts: now, Token: token,
				OrgId: orgID}},
		{fmt.Sprintf("v1/%s/%s/json", orgID, uniqIDTopic), &mqtt.DataPoint{
			Attr: "metadata", MapVal: map[string]string{"aaa": "bbb"}},
			&message.MQTTIngestorOut{UniqId: uniqIDTopic, Attr: "metadata",
				MapVal: map[string]string{"aaa": "bbb"}, Token: token,
				OrgId: orgID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttQueue := queue.NewFake()
			mqttSub, err := mqttQueue.Subscribe("")
			require.NoError(t, err)

			parserQueue := queue.NewFake()
			parserSub, err := parserQueue.Subscribe("")
			require.NoError(t, err)

			ing := Ingestor{
				mqttSub:   mqttSub,
				parserPub: parserQueue,
			}
			go func() {
				ing.parseMessages()
			}()

			var bPayl []byte
			payl := &mqtt.Payload{Points: []*mqtt.DataPoint{lTest.inpPoint},
				Token: token}

			if strings.HasSuffix(lTest.inpTopic, "json") {
				bPayl, err = protojson.Marshal(payl)
			} else {
				bPayl, err = proto.Marshal(payl)
			}
			t.Logf("bPayl: %s", bPayl)
			require.NoError(t, err)

			require.NoError(t, mqttQueue.Publish(lTest.inpTopic, bPayl))

			select {
			case msg := <-parserSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, pubTopic, msg.Topic())

				mIOut := &message.MQTTIngestorOut{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), mIOut))
				t.Logf("mIOut: %#v", mIOut)
				// Normalize generated trace ID.
				lTest.res.TraceId = mIOut.TraceId

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, mIOut) {
					t.Fatalf("Expected, actual: %#v, %#v", lTest.res, mIOut)
				}
			case <-time.After(250 * time.Millisecond):
				t.Error("Message timed out")
			}
		})
	}
}

func TestParseMessagesError(t *testing.T) {
	t.Parallel()

	orgID := uuid.New().String()
	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpPayl  []byte
	}{
		// Bad topic.
		{"v1", nil},
		{fmt.Sprintf("v1/%s/%s/json/aaa", orgID, uniqID), nil},
		{fmt.Sprintf("v2/%s/%s", orgID, uniqID), nil},
		// Bad payload.
		{fmt.Sprintf("v1/%s/%s", orgID, uniqID), []byte("aaa")},
		{fmt.Sprintf("v1/%s/%s/json", orgID, uniqID), []byte("aaa")},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttQueue := queue.NewFake()
			mqttSub, err := mqttQueue.Subscribe("")
			require.NoError(t, err)

			parserQueue := queue.NewFake()
			parserSub, err := parserQueue.Subscribe("")
			require.NoError(t, err)

			ing := Ingestor{
				mqttSub:   mqttSub,
				parserPub: parserQueue,
			}
			go func() {
				ing.parseMessages()
			}()

			require.NoError(t, mqttQueue.Publish(lTest.inpTopic, lTest.inpPayl))

			select {
			case msg := <-parserSub.C():
				t.Errorf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}

func TestDataPointToMIOut(t *testing.T) {
	t.Parallel()

	orgID := uuid.New().String()
	traceID := uuid.New().String()
	token := uuid.New().String()
	uniqIDTopic := random.String(16)
	uniqIDPoint := random.String(16)
	now := timestamppb.Now()

	tests := []struct {
		inpTopicParts []string
		inpPoint      *mqtt.DataPoint
		res           *message.MQTTIngestorOut
	}{
		{[]string{"v1", orgID, uniqIDTopic}, &mqtt.DataPoint{Attr: "motion",
			ValOneof: &mqtt.DataPoint_IntVal{IntVal: 123}},
			&message.MQTTIngestorOut{UniqId: uniqIDTopic, Attr: "motion",
				ValOneof: &message.MQTTIngestorOut_IntVal{IntVal: 123},
				Token:    token, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID}, &mqtt.DataPoint{UniqId: uniqIDPoint,
			Attr: "temp", ValOneof: &mqtt.DataPoint_Fl64Val{Fl64Val: 20.3}},
			&message.MQTTIngestorOut{UniqId: uniqIDPoint, Attr: "temp",
				ValOneof: &message.MQTTIngestorOut_Fl64Val{Fl64Val: 20.3},
				Token:    token, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID, uniqIDTopic}, &mqtt.DataPoint{Attr: "power",
			ValOneof: &mqtt.DataPoint_StrVal{StrVal: "batt"}, Ts: now},
			&message.MQTTIngestorOut{UniqId: uniqIDTopic, Attr: "power",
				ValOneof: &message.MQTTIngestorOut_StrVal{StrVal: "batt"},
				Ts:       now, Token: token, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID}, &mqtt.DataPoint{UniqId: uniqIDPoint,
			Attr: "leak", ValOneof: &mqtt.DataPoint_BoolVal{BoolVal: true}},
			&message.MQTTIngestorOut{UniqId: uniqIDPoint, Attr: "leak",
				ValOneof: &message.MQTTIngestorOut_BoolVal{BoolVal: true},
				Token:    token, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID, uniqIDTopic}, &mqtt.DataPoint{Attr: "metadata",
			MapVal: map[string]string{"aaa": "bbb"}, Ts: now},
			&message.MQTTIngestorOut{UniqId: uniqIDTopic, Attr: "metadata",
				MapVal: map[string]string{"aaa": "bbb"}, Ts: now, Token: token,
				OrgId: orgID, TraceId: traceID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can convert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := dataPointToMIOut(traceID, token, lTest.inpTopicParts,
				lTest.inpPoint)
			t.Logf("res: %#v", res)
			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.res, res) {
				t.Fatalf("Expected, actual: %#v, %#v", lTest.res, res)
			}
		})
	}
}
