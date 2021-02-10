// +build !unit

package test

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
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDecodeMessages(t *testing.T) {
	orgID := uuid.NewString()
	uniqIDPoint := "ing-" + random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	pointToken := uuid.NewString()
	uniqIDTopic := "ing-" + random.String(16)
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
			var bPayl []byte
			var err error
			payl := &mqtt.Payload{Points: lTest.inpPoints,
				Token: lTest.inpPaylToken}

			if lTest.inpTopicParts[len(lTest.inpTopicParts)-1] == "json" {
				bPayl, err = protojson.Marshal(payl)
			} else {
				bPayl, err = proto.Marshal(payl)
			}
			require.NoError(t, err)
			t.Logf("bPayl: %s", bPayl)

			require.NoError(t, globalMQTTQueue.Publish(strings.Join(
				lTest.inpTopicParts, "/"), bPayl))

			for i, res := range lTest.res {
				select {
				case msg := <-globalDecoderSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, globalDecoderPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.Point.TraceId
					// Normalize timestamps.
					if lTest.inpPoints[i].Ts == nil {
						require.WithinDuration(t, time.Now(),
							vIn.Point.Ts.AsTime(), 5*time.Second)
						res.Point.Ts = vIn.Point.Ts
					}

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, vIn)
					}
				case <-time.After(5 * time.Second):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestDecodeMessagesError(t *testing.T) {
	orgID := uuid.NewString()
	uniqID := "ing-" + random.String(16)

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

		t.Run(fmt.Sprintf("Cannot decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				lTest.inpPayl))

			select {
			case msg := <-globalDecoderSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
