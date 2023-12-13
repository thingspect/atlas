//go:build !unit

package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/common"
	"github.com/thingspect/proto/go/mqtt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 6 * time.Second

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
		{
			[]string{"v1", orgID, "json"}, "", []*common.DataPoint{
				{
					UniqId: uniqIDPoint, Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
					Token: pointToken,
				}, {
					UniqId: uniqIDPoint, Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 321}, Ts: now,
					Token: pointToken,
				},
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqIDPoint,
						Attr:   "count", ValOneof: &common.DataPoint_IntVal{
							IntVal: 123,
						}, Ts: now, Token: pointToken,
					}, OrgId: orgID,
				}, {
					Point: &common.DataPoint{
						UniqId: uniqIDPoint,
						Attr:   "count", ValOneof: &common.DataPoint_IntVal{
							IntVal: 321,
						}, Ts: now, Token: pointToken,
					}, OrgId: orgID,
				},
			},
		},
		{
			[]string{"v1", orgID, uniqIDTopic}, paylToken, []*common.DataPoint{
				{
					Attr:     "temp_c",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				},
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqIDTopic, Attr: "temp_c",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
						Token:    paylToken,
					}, OrgId: orgID,
				},
			},
		},
		{
			[]string{"v1", orgID, uniqIDTopic, "json"},
			paylToken,
			[]*common.DataPoint{
				{
					Attr:     "power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
				},
			},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqIDTopic, Attr: "power",
						ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
						Token:    paylToken,
					}, OrgId: orgID,
				},
			},
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			var bPayl []byte
			var err error
			payl := &mqtt.Payload{
				Points: lTest.inpPoints, Token: lTest.inpPaylToken,
			}

			if lTest.inpTopicParts[len(lTest.inpTopicParts)-1] == "json" {
				bPayl, err = protojson.Marshal(payl)
			} else {
				bPayl, err = proto.Marshal(payl)
			}
			require.NoError(t, err)
			t.Logf("bPayl: %s", bPayl)

			require.NoError(t, globalMQTTQueue.Publish(strings.Join(
				lTest.inpTopicParts, "/"), bPayl))

			// Don't stop the flow of execution (assert) to avoid leaving
			// messages orphaned in the queue.
			for i, res := range lTest.res {
				select {
				//nolint:testifylint // above
				case msg := <-globalVInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					assert.Equal(t, globalVInPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					assert.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.GetPoint().GetTraceId()
					// Normalize timestamp.
					if lTest.inpPoints[i].GetTs() == nil {
						assert.WithinDuration(t, time.Now(),
							vIn.GetPoint().GetTs().AsTime(), testTimeout)
						res.Point.Ts = vIn.GetPoint().GetTs()
					}

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Errorf("\nExpect: %+v\nActual: %+v", res, vIn)
					}
				case <-time.After(testTimeout):
					t.Error("Message timed out")
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
			case msg := <-globalVInSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
