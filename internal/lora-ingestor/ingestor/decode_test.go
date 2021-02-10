// +build !integration

package ingestor

import (
	"encoding/base64"
	"fmt"
	"io"
	"testing"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDecodeGateways(t *testing.T) {
	t.Parallel()

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

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttGWQueue := queue.NewFake()
			mqttGWSub, err := mqttGWQueue.Subscribe("")
			require.NoError(t, err)

			decoderQueue := queue.NewFake()
			decoderSub, err := decoderQueue.Subscribe("")
			require.NoError(t, err)
			decoderPubGWTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttGWSub: mqttGWSub,

				decoderQueue:      decoderQueue,
				decoderPubGWTopic: decoderPubGWTopic,
			}
			go func() {
				ing.decodeGateways()
			}()

			bInpProto, err := proto.Marshal(lTest.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, mqttGWQueue.Publish(lTest.inpTopic, bInpProto))

			for _, res := range lTest.res {
				select {
				case msg := <-decoderSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, decoderPubGWTopic, msg.Topic())

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
				case <-time.After(2 * time.Second):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestDecodeGatewaysError(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpPayl  []byte
	}{
		// Bad topic.
		{"lora", nil},
		// Bad payload.
		{"lora/gateway/" + uniqID + "/event/up", []byte("ing-aaa")},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttGWQueue := queue.NewFake()
			mqttGWSub, err := mqttGWQueue.Subscribe("")
			require.NoError(t, err)

			decoderQueue := queue.NewFake()
			decoderSub, err := decoderQueue.Subscribe("")
			require.NoError(t, err)
			decoderPubGWTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttGWSub: mqttGWSub,

				decoderQueue:      decoderQueue,
				decoderPubGWTopic: decoderPubGWTopic,
			}
			go func() {
				ing.decodeGateways()
			}()

			require.NoError(t, mqttGWQueue.Publish(lTest.inpTopic,
				lTest.inpPayl))

			select {
			case msg := <-decoderSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}

func TestDecodeDevices(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)

	bData := random.Bytes(10)
	b64Data := base64.StdEncoding.EncodeToString(bData)
	t.Logf("b64Data: %v", b64Data)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		res      []*message.ValidatorIn
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
		}},
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
		}},
		{"lora/application/3/device/" + uniqID + "/event/ack", &as.AckEvent{
			Acknowledged: true}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"acknowledged":true}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_StrVal{StrVal: "OK"}},
				SkipToken: true},
		}},
		{"lora/application/4/device/" + uniqID + "/event/error", &as.ErrorEvent{
			Type: as.ErrorType_OTAA}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
				ValOneof: &common.DataPoint_StrVal{
					StrVal: `{"type":"OTAA"}`}}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: uniqID, Attr: "error_type",
				ValOneof: &common.DataPoint_StrVal{StrVal: "OTAA"}},
				SkipToken: true},
		}},
		{"lora/application/5/device/" + uniqID + "/event/txack",
			&as.TxAckEvent{}, []*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "raw_device",
					ValOneof: &common.DataPoint_StrVal{StrVal: `{}`}},
					SkipToken: true},
				{Point: &common.DataPoint{UniqId: uniqID,
					Attr: "ack_gateway_tx", ValOneof: &common.DataPoint_BoolVal{
						BoolVal: true}}, SkipToken: true},
			}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttDevQueue := queue.NewFake()
			mqttDevSub, err := mqttDevQueue.Subscribe("")
			require.NoError(t, err)

			decoderQueue := queue.NewFake()
			decoderSub, err := decoderQueue.Subscribe("")
			require.NoError(t, err)
			decoderPubDevTopic := "topic-" + random.String(10)
			decoderPubDataTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttDevSub: mqttDevSub,

				decoderQueue:        decoderQueue,
				decoderPubDevTopic:  decoderPubDevTopic,
				decoderPubDataTopic: decoderPubDataTopic,
			}
			go func() {
				ing.decodeDevices()
			}()

			bInpProto, err := proto.Marshal(lTest.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, mqttDevQueue.Publish(lTest.inpTopic, bInpProto))

			for _, res := range lTest.res {
				select {
				case msg := <-decoderSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, decoderPubDevTopic, msg.Topic())

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
				case <-time.After(2 * time.Second):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestDecodeDevicesError(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpPayl  []byte
	}{
		// Bad topic.
		{"lora", nil},
		// Bad payload.
		{"lora/application/1/device/" + uniqID + "/event/up",
			[]byte("ing-aaa")},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			mqttDevQueue := queue.NewFake()
			mqttDevSub, err := mqttDevQueue.Subscribe("")
			require.NoError(t, err)

			decoderQueue := queue.NewFake()
			decoderSub, err := decoderQueue.Subscribe("")
			require.NoError(t, err)
			decoderPubDevTopic := "topic-" + random.String(10)
			decoderPubDataTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttDevSub: mqttDevSub,

				decoderQueue:        decoderQueue,
				decoderPubDevTopic:  decoderPubDevTopic,
				decoderPubDataTopic: decoderPubDataTopic,
			}
			go func() {
				ing.decodeDevices()
			}()

			require.NoError(t, mqttDevQueue.Publish(lTest.inpTopic,
				lTest.inpPayl))

			select {
			case msg := <-decoderSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}

func TestPointToVIn(t *testing.T) {
	t.Parallel()

	traceID := uuid.NewString()
	uniqID := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	tests := []struct {
		inp *decode.Point
		res *message.ValidatorIn
	}{
		{&decode.Point{Attr: "data_rate", Value: int32(3)}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "data_rate",
				ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "data_rate", Value: 3}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "data_rate",
				ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "data_rate", Value: int64(3)}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "data_rate",
				ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "snr", Value: 7.8}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "snr",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 7.8}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "snr", Value: float32(7.0)}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "snr",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 7.0}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "ack", Value: "OK"}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_StrVal{StrVal: "OK"}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "adr", Value: false}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "adr",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: false}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "ack", Value: []byte{0x00}}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_BytesVal{BytesVal: []byte{0x00}},
				Ts:       now, TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "error", Value: io.EOF}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "error", Ts: now,
				TraceId: traceID}, SkipToken: true}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can convert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := pointToVIn(traceID, uniqID, lTest.inp, now)
			t.Logf("res: %+v", res)

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.res, res) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, res)
			}
		})
	}
}
