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

func TestDecodeMessages(t *testing.T) {
	t.Parallel()

	orgID := uuid.NewString()
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	pointToken := uuid.NewString()
	uniqIDTopic := random.String(16)
	paylToken := uuid.NewString()

	tests := []struct {
		inpTopicParts []string
		inpPaylToken  string
		inpPoints     []*common.DataPoint
		res           []*message.ValidatorIn
	}{
		{[]string{"v1", orgID, "json"}, "", []*common.DataPoint{
			{UniqId: uniqIDPoint, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
				Token: pointToken},
			{UniqId: uniqIDPoint, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 321}, Ts: now,
				Token: pointToken},
		}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqIDPoint,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: now, Token: pointToken}, OrgId: orgID},
			{Point: &common.DataPoint{UniqId: uniqIDPoint,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 321},
				Ts: now, Token: pointToken}, OrgId: orgID},
		}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken, []*common.DataPoint{
			{Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}},
		}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqIDTopic, Attr: "temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Token:    paylToken}, OrgId: orgID},
		}},
		{[]string{"v1", orgID, uniqIDTopic, "json"}, paylToken,
			[]*common.DataPoint{
				{Attr: "power", ValOneof: &common.DataPoint_StrVal{
					StrVal: "batt"}},
			}, []*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqIDTopic, Attr: "power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
					Token:    paylToken}, OrgId: orgID},
			}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttQueue := queue.NewFake()
			mqttSub, err := mqttQueue.Subscribe("")
			require.NoError(t, err)

			ingQueue := queue.NewFake()
			vInSub, err := ingQueue.Subscribe("")
			require.NoError(t, err)
			vInPubTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttSub: mqttSub,

				ingQueue:    ingQueue,
				vInPubTopic: vInPubTopic,
			}
			go func() {
				ing.decodeMessages()
			}()

			var bPayl []byte
			payl := &mqtt.Payload{Points: lTest.inpPoints,
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

			for i, res := range lTest.res {
				select {
				case msg := <-vInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, vInPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.Point.TraceId
					// Normalize timestamps.
					if lTest.inpPoints[i].Ts == nil {
						require.WithinDuration(t, time.Now(),
							vIn.Point.Ts.AsTime(), 2*time.Second)
						res.Point.Ts = vIn.Point.Ts
					}

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, vIn)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestDecodeMessagesError(t *testing.T) {
	t.Parallel()

	orgID := uuid.NewString()
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

		t.Run(fmt.Sprintf("Cannot decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttQueue := queue.NewFake()
			mqttSub, err := mqttQueue.Subscribe("")
			require.NoError(t, err)

			ingQueue := queue.NewFake()
			vInSub, err := ingQueue.Subscribe("")
			require.NoError(t, err)

			ing := Ingestor{
				mqttSub: mqttSub,

				ingQueue:    ingQueue,
				vInPubTopic: "topic-" + random.String(10),
			}
			go func() {
				ing.decodeMessages()
			}()

			require.NoError(t, mqttQueue.Publish(strings.Join(
				lTest.inpTopicParts, "/"), lTest.inpPayl))

			select {
			case msg := <-vInSub.C():
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

	orgID := uuid.NewString()
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	pointToken := uuid.NewString()
	traceID := uuid.NewString()
	uniqIDTopic := random.String(16)
	paylToken := uuid.NewString()

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
			t.Logf("res: %+v", res)

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
