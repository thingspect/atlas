// +build !integration

package ingestor

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
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
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	pointToken := uuid.New().String()
	uniqIDTopic := random.String(16)
	paylToken := uuid.New().String()

	tests := []struct {
		inpTopicParts []string
		inpPaylToken  string
		inpPoint      *common.DataPoint
		res           *message.ValidatorIn
	}{
		{[]string{"v1", orgID, "json"}, "",
			&common.DataPoint{UniqId: uniqIDPoint, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
				Token: pointToken},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDPoint,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: now, Token: pointToken}, OrgId: orgID}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&common.DataPoint{Attr: "temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDTopic,
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Token: paylToken}, OrgId: orgID}},
		{[]string{"v1", orgID, uniqIDTopic, "json"}, paylToken,
			&common.DataPoint{Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"}},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDTopic,
				Attr:     "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
				Token:    paylToken}, OrgId: orgID}},
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
			payl := &mqtt.Payload{Points: []*common.DataPoint{lTest.inpPoint},
				Token: lTest.inpPaylToken}

			if lTest.inpTopicParts[len(lTest.inpTopicParts)-1] == "json" {
				bPayl, err = protojson.Marshal(payl)
			} else {
				bPayl, err = proto.Marshal(payl)
			}
			require.NoError(t, err)
			t.Logf("bPayl: %s", bPayl)

			require.NoError(t, mqttQueue.Publish(strings.Join(
				lTest.inpTopicParts, "/"), bPayl))

			select {
			case msg := <-parserSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, parserPubTopic, msg.Topic())

				vIn := &message.ValidatorIn{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
				t.Logf("vIn: %+v", vIn)

				// Normalize generated trace ID.
				lTest.res.Point.TraceId = vIn.Point.TraceId
				// Normalize timestamps.
				if lTest.inpPoint.Ts == nil {
					require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
						2*time.Second)
					lTest.res.Point.Ts = vIn.Point.Ts
				}

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vIn) {
					t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, vIn)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Message timed out")
			}
		})
	}
}

func TestParseMessagesError(t *testing.T) {
	t.Parallel()

	orgID := uuid.New().String()
	uniqID := random.String(16)

	tests := []struct {
		inpTopicParts []string
		inpPayl       []byte
	}{
		// Bad topic.
		{[]string{"v1"}, nil},
		{[]string{"v1", orgID, uniqID, "json/ing-aaa"}, nil},
		{[]string{"v2", orgID, uniqID}, nil},
		// Bad payload.
		{[]string{"v1", orgID, uniqID}, []byte("ing-aaa")},
		{[]string{"v1", orgID, uniqID, "json"}, []byte("ing-aaa")},
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

			require.NoError(t, mqttQueue.Publish(strings.Join(
				lTest.inpTopicParts, "/"), lTest.inpPayl))

			select {
			case msg := <-parserSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
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
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	pointToken := uuid.New().String()
	traceID := uuid.New().String()
	uniqIDTopic := random.String(16)
	paylToken := uuid.New().String()

	tests := []struct {
		inpTopicParts []string
		inpPaylToken  string
		inpPoint      *common.DataPoint
		res           *message.ValidatorIn
	}{
		{[]string{"v1", orgID}, "",
			&common.DataPoint{UniqId: uniqIDPoint, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
				Token: pointToken},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDPoint,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: now, Token: pointToken, TraceId: traceID}, OrgId: orgID}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&common.DataPoint{Attr: "temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDTopic,
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Token: paylToken, TraceId: traceID}, OrgId: orgID}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&common.DataPoint{Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"}},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDTopic,
				Attr:     "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
				Token:    paylToken, TraceId: traceID}, OrgId: orgID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can convert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			// Save original TS before lTest.inpPoint is modified in-place.
			origTS := lTest.inpPoint.Ts

			res := dataPointToVIn(traceID, lTest.inpPaylToken,
				lTest.inpTopicParts, lTest.inpPoint)

			// Normalize timestamps.
			if origTS == nil {
				require.WithinDuration(t, time.Now(), res.Point.Ts.AsTime(),
					2*time.Second)
				lTest.res.Point.Ts = res.Point.Ts
			}

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.res, res) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, res)
			}
		})
	}
}
