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

func TestParseMessages(t *testing.T) {
	orgID := uuid.New().String()
	paylToken := uuid.New().String()
	pointToken := uuid.New().String()
	uniqIDTopic := random.String(16)
	uniqIDPoint := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	tests := []struct {
		inpTopicParts []string
		inpPaylToken  string
		inpPoint      *common.DataPoint
		res           *message.ValidatorIn
	}{
		{[]string{"v1", orgID, "json"}, "",
			&common.DataPoint{UniqId: uniqIDPoint, Attr: "ing-motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
				Token: pointToken},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDPoint,
				Attr:     "ing-motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
				Token: pointToken}, OrgId: orgID}},
		{[]string{"v1", orgID, uniqIDTopic}, paylToken,
			&common.DataPoint{Attr: "ing-temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDTopic,
				Attr:     "ing-temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Token:    paylToken}, OrgId: orgID}},
		{[]string{"v1", orgID, uniqIDTopic, "json"}, paylToken,
			&common.DataPoint{Attr: "ing-power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"}},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqIDTopic,
				Attr:     "ing-power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
				Token:    paylToken}, OrgId: orgID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			var bPayl []byte
			var err error
			payl := &mqtt.Payload{Points: []*common.DataPoint{lTest.inpPoint},
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

			select {
			case msg := <-globalParserSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, globalParserPubTopic, msg.Topic())

				vIn := &message.ValidatorIn{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
				t.Logf("vIn: %#v", vIn)

				// Normalize generated trace ID.
				lTest.res.Point.TraceId = vIn.Point.TraceId
				// Normalize timestamps.
				if lTest.inpPoint.Ts == nil {
					require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
						5*time.Second)
					lTest.res.Point.Ts = vIn.Point.Ts
				}

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vIn) {
					t.Fatalf("Expected, actual: %#v, %#v", lTest.res, vIn)
				}
			case <-time.After(5 * time.Second):
				t.Error("Message timed out")
			}
		})
	}
}

func TestParseMessagesError(t *testing.T) {
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

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				lTest.inpPayl))

			select {
			case msg := <-globalParserSub.C():
				t.Errorf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
