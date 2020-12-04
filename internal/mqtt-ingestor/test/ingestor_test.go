// +build !unit

package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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
	now := timestamppb.Now()

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
				MapVal: map[string]string{"ing-aaa": "ing-bbb"}, Token: paylToken,
				OrgId: orgID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			var bPayl []byte
			var err error
			payl := &mqtt.Payload{Points: []*mqtt.DataPoint{lTest.inpPoint},
				Token: lTest.inpPaylToken}

			if strings.HasSuffix(lTest.inpTopic, "json") {
				bPayl, err = protojson.Marshal(payl)
			} else {
				bPayl, err = proto.Marshal(payl)
			}
			require.NoError(t, err)
			t.Logf("bPayl: %s", bPayl)

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic, bPayl))

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
				lTest.res.TraceId = vIn.TraceId
				// Normalize timestamps.
				if lTest.inpPoint.Ts == nil {
					require.WithinDuration(t, time.Now(), vIn.Ts.AsTime(),
						5*time.Second)
					lTest.res.Ts = vIn.Ts
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
