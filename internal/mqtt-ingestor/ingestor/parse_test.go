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
	paylToken := uuid.New().String()
	pointToken := uuid.New().String()
	uniqIDTopic := random.String(16)
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	tests := []struct {
		inpTopic     string
		inpPaylToken string
		inpPoint     *mqtt.DataPoint
		res          *message.ValidatorIn
	}{
		{fmt.Sprintf("v1/%s/%s", orgID, uniqIDTopic), paylToken,
			&mqtt.DataPoint{Attr: "motion",
				ValOneof: &mqtt.DataPoint_IntVal{IntVal: 123}},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "motion",
				ValOneof: &message.ValidatorIn_IntVal{IntVal: 123},
				Token:    paylToken, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/json", orgID), "",
			&mqtt.DataPoint{UniqId: uniqIDPoint, Attr: "temp",
				ValOneof: &mqtt.DataPoint_Fl64Val{Fl64Val: 20.3},
				Token:    pointToken},
			&message.ValidatorIn{UniqId: uniqIDPoint, Attr: "temp",
				ValOneof: &message.ValidatorIn_Fl64Val{Fl64Val: 20.3},
				Token:    pointToken, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/%s/json", orgID, uniqIDTopic), paylToken,
			&mqtt.DataPoint{Attr: "power",
				ValOneof: &mqtt.DataPoint_StrVal{StrVal: "batt"}, Ts: now},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "power",
				ValOneof: &message.ValidatorIn_StrVal{StrVal: "batt"}, Ts: now,
				Token: paylToken, OrgId: orgID}},
		{fmt.Sprintf("v1/%s", orgID), paylToken,
			&mqtt.DataPoint{UniqId: uniqIDPoint, Attr: "leak",
				ValOneof: &mqtt.DataPoint_BoolVal{BoolVal: true}},
			&message.ValidatorIn{UniqId: uniqIDPoint, Attr: "leak",
				ValOneof: &message.ValidatorIn_BoolVal{BoolVal: true},
				Token:    paylToken, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/json", orgID), paylToken,
			&mqtt.DataPoint{UniqId: uniqIDPoint, Attr: "raw",
				ValOneof: &mqtt.DataPoint_BytesVal{BytesVal: []byte{0x00}}},
			&message.ValidatorIn{UniqId: uniqIDPoint, Attr: "raw",
				ValOneof: &message.ValidatorIn_BytesVal{BytesVal: []byte{0x00}},
				Token:    paylToken, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/%s", orgID, uniqIDTopic), paylToken,
			&mqtt.DataPoint{Attr: "metadata",
				MapVal: map[string]string{"ing-aaa": "ing-bbb"}, Ts: now},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "metadata",
				MapVal: map[string]string{"ing-aaa": "ing-bbb"}, Ts: now,
				Token: paylToken, OrgId: orgID}},
		{fmt.Sprintf("v1/%s/%s/json", orgID, uniqIDTopic), paylToken,
			&mqtt.DataPoint{Attr: "metadata",
				MapVal: map[string]string{"ing-aaa": "ing-bbb"}},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "metadata",
				MapVal: map[string]string{"ing-aaa": "ing-bbb"},
				Token:  paylToken, OrgId: orgID}},
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
			parserPubTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttSub:        mqttSub,
				parserQueue:    parserQueue,
				parserPubTopic: parserPubTopic,
			}
			go func() {
				ing.parseMessages()
			}()

			var bPayl []byte
			payl := &mqtt.Payload{Points: []*mqtt.DataPoint{lTest.inpPoint},
				Token: lTest.inpPaylToken}

			if strings.HasSuffix(lTest.inpTopic, "json") {
				bPayl, err = protojson.Marshal(payl)
			} else {
				bPayl, err = proto.Marshal(payl)
			}
			require.NoError(t, err)
			t.Logf("bPayl: %s", bPayl)

			require.NoError(t, mqttQueue.Publish(lTest.inpTopic, bPayl))

			select {
			case msg := <-parserSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, parserPubTopic, msg.Topic())

				vIn := &message.ValidatorIn{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
				t.Logf("vIn: %#v", vIn)

				// Normalize generated trace ID.
				lTest.res.TraceId = vIn.TraceId
				// Normalize timestamps.
				if lTest.inpPoint.Ts == nil {
					require.WithinDuration(t, time.Now(), vIn.Ts.AsTime(),
						2*time.Second)
					lTest.res.Ts = vIn.Ts
				}

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vIn) {
					t.Fatalf("Expected, actual: %#v, %#v", lTest.res, vIn)
				}
			case <-time.After(2 * time.Second):
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
		{fmt.Sprintf("v1/%s/%s/json/ing-aaa", orgID, uniqID), nil},
		{fmt.Sprintf("v2/%s/%s", orgID, uniqID), nil},
		// Bad payload.
		{fmt.Sprintf("v1/%s/%s", orgID, uniqID), []byte("ing-aaa")},
		{fmt.Sprintf("v1/%s/%s/json", orgID, uniqID), []byte("ing-aaa")},
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
				mqttSub:        mqttSub,
				parserQueue:    parserQueue,
				parserPubTopic: "topic-" + random.String(10),
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

func TestDataPointToVIn(t *testing.T) {
	t.Parallel()

	orgID := uuid.New().String()
	traceID := uuid.New().String()
	paylToken := uuid.New().String()
	pointToken := uuid.New().String()
	uniqIDTopic := random.String(16)
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	tests := []struct {
		inpTopicParts []string
		inpPaylToken  string
		inpPoint      *mqtt.DataPoint
		res           *message.ValidatorIn
	}{
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&mqtt.DataPoint{Attr: "motion",
				ValOneof: &mqtt.DataPoint_IntVal{IntVal: 123}},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "motion",
				ValOneof: &message.ValidatorIn_IntVal{IntVal: 123},
				Token:    paylToken, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID}, "",
			&mqtt.DataPoint{UniqId: uniqIDPoint, Attr: "temp",
				ValOneof: &mqtt.DataPoint_Fl64Val{Fl64Val: 20.3},
				Token:    pointToken},
			&message.ValidatorIn{UniqId: uniqIDPoint, Attr: "temp",
				ValOneof: &message.ValidatorIn_Fl64Val{Fl64Val: 20.3},
				Token:    pointToken, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&mqtt.DataPoint{Attr: "power",
				ValOneof: &mqtt.DataPoint_StrVal{StrVal: "batt"}, Ts: now},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "power",
				ValOneof: &message.ValidatorIn_StrVal{StrVal: "batt"},
				Ts:       now, Token: paylToken, OrgId: orgID,
				TraceId: traceID}},
		{[]string{"v1", orgID}, paylToken,
			&mqtt.DataPoint{UniqId: uniqIDPoint, Attr: "leak",
				ValOneof: &mqtt.DataPoint_BoolVal{BoolVal: true}},
			&message.ValidatorIn{UniqId: uniqIDPoint, Attr: "leak",
				ValOneof: &message.ValidatorIn_BoolVal{BoolVal: true},
				Token:    paylToken, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID}, paylToken,
			&mqtt.DataPoint{UniqId: uniqIDPoint, Attr: "raw",
				ValOneof: &mqtt.DataPoint_BytesVal{BytesVal: []byte{0x00}}},
			&message.ValidatorIn{UniqId: uniqIDPoint, Attr: "raw",
				ValOneof: &message.ValidatorIn_BytesVal{BytesVal: []byte{0x00}},
				Token:    paylToken, OrgId: orgID, TraceId: traceID}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&mqtt.DataPoint{Attr: "metadata",
				MapVal: map[string]string{"ing-aaa": "ing-bbb"}, Ts: now},
			&message.ValidatorIn{UniqId: uniqIDTopic, Attr: "metadata",
				MapVal: map[string]string{"ing-aaa": "ing-bbb"}, Ts: now,
				Token: paylToken, OrgId: orgID, TraceId: traceID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can convert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := dataPointToVIn(traceID, lTest.inpPaylToken,
				lTest.inpTopicParts, lTest.inpPoint)
			t.Logf("res: %#v", res)

			// Normalize timestamps.
			if lTest.inpPoint.Ts == nil {
				require.WithinDuration(t, time.Now(), res.Ts.AsTime(),
					2*time.Second)
				lTest.res.Ts = res.Ts
			}

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.res, res) {
				t.Fatalf("Expected, actual: %#v, %#v", lTest.res, res)
			}
		})
	}
}
