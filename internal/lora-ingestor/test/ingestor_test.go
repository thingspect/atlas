// +build !unit

package test

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestParseGateways(t *testing.T) {
	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		res      []*message.ValidatorIn
	}{
		{"lora/gateway/" + uniqID + "/event/up", &gw.UplinkFrame{
			RxInfo: &gw.UplinkRXInfo{Rssi: -74}}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_gateway",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"rxInfo":{"rssi":-74}}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "lora_rssi",
				ValOneof: &common.DataPoint_IntVal{IntVal: -74}},
				SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "channel",
				ValOneof: &common.DataPoint_IntVal{IntVal: 0}},
				SkipToken: true},
		}},
		{"lora/gateway/" + uniqID + "/event/stats", &gw.GatewayStats{
			RxPacketsReceivedOk: 2}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_gateway",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"rxPacketsReceivedOK":2}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "rx_received_valid",
				ValOneof: &common.DataPoint_IntVal{IntVal: 2}},
				SkipToken: true},
		}},
		{"lora/gateway/" + uniqID + "/event/ack", &gw.DownlinkTXAck{
			Items: []*gw.DownlinkTXAckItem{{Status: gw.TxAckStatus_OK}}},
			[]*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_gateway",
					ValOneof: &common.DataPoint_StrVal{StrVal: `{"items":[{` +
						`"status":"OK"}]}`}}, SkipToken: true},
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
					ValOneof: &common.DataPoint_StrVal{StrVal: "OK"}},
					SkipToken: true},
			}},
		{"lora/gateway/" + uniqID + "/event/exec",
			&gw.GatewayCommandExecResponse{Stdout: []byte("STDOUT")},
			[]*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_gateway",
					ValOneof: &common.DataPoint_StrVal{
						StrVal: `{"stdout":"U1RET1VU"}`}}, SkipToken: true},
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "exec_stdout",
					ValOneof: &common.DataPoint_StrVal{StrVal: "STDOUT"}},
					SkipToken: true},
			}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			bInpProto, err := proto.Marshal(lTest.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				bInpProto))

			for _, res := range lTest.res {
				select {
				case msg := <-globalParserGWSub.C():
					msg.Ack()
					t.Logf("GW msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, globalParserPubGWTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.Point.TraceId
					// Normalize timestamps.
					require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
						2*time.Second)
					res.Point.Ts = vIn.Point.Ts

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

func TestParseDevices(t *testing.T) {
	uniqID := random.String(16)

	bData := random.Bytes(10)
	b64Data := base64.StdEncoding.EncodeToString(bData)
	t.Logf("b64Data: %v", b64Data)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		resVIn   []*message.ValidatorIn
		resPIn   *message.ParserIn
	}{
		{"lora/application/1/device/" + uniqID + "/event/up", &as.UplinkEvent{
			Data: bData}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
				ValOneof: &common.DataPoint_StrVal{StrVal: fmt.Sprintf(
					`{"data":"%s"}`, b64Data)}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "adr",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: false}},
				SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "data_rate",
				ValOneof: &common.DataPoint_IntVal{IntVal: 0}},
				SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "confirmed",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: false}},
				SkipToken: true},
		}, &message.ParserIn{UniqId: uniqID, Data: bData}},
		{"lora/application/2/device/" + uniqID + "/event/join", &as.JoinEvent{
			RxInfo: []*gw.UplinkRXInfo{{}}, Dr: 3}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"rxInfo":[{}],"dr":3}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "join",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: true}},
				SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "channel",
				ValOneof: &common.DataPoint_IntVal{IntVal: 0}},
				SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "data_rate",
				ValOneof: &common.DataPoint_IntVal{IntVal: 3}},
				SkipToken: true},
		}, nil},
		{"lora/application/3/device/" + uniqID + "/event/ack", &as.AckEvent{
			Acknowledged: true}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"acknowledged":true}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_StrVal{StrVal: "OK"}},
				SkipToken: true},
		}, nil},
		{"lora/application/4/device/" + uniqID + "/event/error", &as.ErrorEvent{
			Type: as.ErrorType_OTAA}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"type":"OTAA"}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "error_type",
				ValOneof: &common.DataPoint_StrVal{StrVal: "OTAA"}},
				SkipToken: true},
		}, nil},
		{"lora/application/5/device/" + uniqID + "/event/txack",
			&as.TxAckEvent{}, []*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
					ValOneof: &common.DataPoint_StrVal{StrVal: `{}`}},
					SkipToken: true},
				{Point: &common.DataPoint{UniqId: uniqID,
					Attr: "ack_gateway_tx", ValOneof: &common.DataPoint_BoolVal{
						BoolVal: true}}, SkipToken: true},
			}, nil},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			bInpProto, err := proto.Marshal(lTest.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				bInpProto))

			for _, res := range lTest.resVIn {
				select {
				case msg := <-globalParserDevSub.C():
					msg.Ack()
					t.Logf("Dev msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, globalParserPubDevTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.Point.TraceId
					// Normalize timestamps.
					require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
						2*time.Second)
					res.Point.Ts = vIn.Point.Ts

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.resVIn,
							vIn)
					}
				case <-time.After(5 * time.Second):
					t.Fatal("Message timed out")
				}
			}

			if lTest.resPIn != nil {
				select {
				case msg := <-globalParserDataSub.C():
					msg.Ack()
					t.Logf("Data msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, globalParserPubDataTopic, msg.Topic())

					pIn := &message.ParserIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), pIn))
					t.Logf("pIn: %+v", pIn)

					// Normalize generated trace ID.
					lTest.resPIn.TraceId = pIn.TraceId
					// Normalize timestamps.
					require.WithinDuration(t, time.Now(), pIn.Ts.AsTime(),
						2*time.Second)
					lTest.resPIn.Ts = pIn.Ts

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(lTest.resPIn, pIn) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.resPIn,
							pIn)
					}
				case <-time.After(5 * time.Second):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestParseGatewaysDevicesError(t *testing.T) {
	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpPayl  []byte
	}{
		// Bad topic.
		{"lora", nil},
		// Bad payload.
		{"lora/gateway/" + uniqID + "/event/up", []byte("ing-aaa")},
		{"lora/application/1/device/" + uniqID + "/event/up",
			[]byte("ing-aaa")},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				lTest.inpPayl))

			select {
			case msg := <-globalParserGWSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case msg := <-globalParserDevSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
