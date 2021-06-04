// +build !unit

package test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 6 * time.Second

func TestDecodeGateways(t *testing.T) {
	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		res      []*message.ValidatorIn
	}{
		{
			"lora/gateway/" + uniqID + "/event/up", &gw.UplinkFrame{
				RxInfo: &gw.UplinkRXInfo{Rssi: -74},
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_gateway",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"rxInfo":{"rssi":-74}}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "lora_rssi",
						ValOneof: &common.DataPoint_IntVal{IntVal: -74},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "channel",
						ValOneof: &common.DataPoint_IntVal{IntVal: 0},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/gateway/" + uniqID + "/event/stats", &gw.GatewayStats{
				RxPacketsReceivedOk: 2,
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_gateway",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"rxPacketsReceivedOK":2}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "rx_received_valid",
						ValOneof: &common.DataPoint_IntVal{IntVal: 2},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/gateway/" + uniqID + "/event/ack", &gw.DownlinkTXAck{
				Items: []*gw.DownlinkTXAckItem{{Status: gw.TxAckStatus_OK}},
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_gateway",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"items":[{"status":"OK"}]}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "ack",
						ValOneof: &common.DataPoint_StrVal{StrVal: "OK"},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/gateway/" + uniqID + "/event/exec",
			&gw.GatewayCommandExecResponse{Stdout: []byte("STDOUT")},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_gateway",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"stdout":"U1RET1VU"}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "exec_stdout",
						ValOneof: &common.DataPoint_StrVal{StrVal: "STDOUT"},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/gateway/" + uniqID + "/state/conn", &gw.ConnState{
				State: gw.ConnState_ONLINE,
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_gateway",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"state":"ONLINE"}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "conn",
						ValOneof: &common.DataPoint_StrVal{StrVal: "ONLINE"},
					}, SkipToken: true,
				},
			},
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			bInpProto, err := proto.Marshal(lTest.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				bInpProto))

			// Don't stop the flow of execution (assert) to avoid leaving
			// messages orphaned in the queue.
			for _, res := range lTest.res {
				select {
				case msg := <-globalVInGWSub.C():
					msg.Ack()
					t.Logf("GW msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					assert.Equal(t, globalVInGWPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					assert.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.Point.TraceId
					// Normalize timestamps.
					assert.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
						testTimeout)
					res.Point.Ts = vIn.Point.Ts

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Errorf("\nExpect: %+v\nActual: %+v", lTest.res, vIn)
					}
				case <-time.After(testTimeout):
					t.Error("Message timed out")
				}
			}
		})
	}
}

func TestDecodeDevices(t *testing.T) {
	uniqID := random.String(16)

	bData := random.Bytes(10)
	b64Data := base64.StdEncoding.EncodeToString(bData)
	t.Logf("b64Data: %v", b64Data)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		resVIn   []*message.ValidatorIn
		resPIn   *message.DecoderIn
	}{
		{
			"lora/application/1/device/" + uniqID + "/event/up",
			&as.UplinkEvent{Data: bData},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: fmt.Sprintf(`{"data":"%s"}`, b64Data),
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_data",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: hex.EncodeToString(bData),
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "adr",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "data_rate",
						ValOneof: &common.DataPoint_IntVal{IntVal: 0},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "confirmed",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
					}, SkipToken: true,
				},
			},
			&message.DecoderIn{UniqId: uniqID, Data: bData},
		},
		{
			"lora/application/2/device/" + uniqID + "/event/join",
			&as.JoinEvent{RxInfo: []*gw.UplinkRXInfo{{}}, Dr: 3},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"rxInfo":[{}],"dr":3}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "join",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "channel",
						ValOneof: &common.DataPoint_IntVal{IntVal: 0},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "data_rate",
						ValOneof: &common.DataPoint_IntVal{IntVal: 3},
					}, SkipToken: true,
				},
			},
			nil,
		},
		{
			"lora/application/3/device/" + uniqID + "/event/ack", &as.AckEvent{
				Acknowledged: true,
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"acknowledged":true}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "ack",
						ValOneof: &common.DataPoint_StrVal{StrVal: "OK"},
					}, SkipToken: true,
				},
			}, nil,
		},
		{
			"lora/application/4/device/" + uniqID + "/event/error",
			&as.ErrorEvent{Type: as.ErrorType_OTAA},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"type":"OTAA"}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "error_type",
						ValOneof: &common.DataPoint_StrVal{StrVal: "OTAA"},
					}, SkipToken: true,
				},
			},
			nil,
		},
		{
			"lora/application/5/device/" + uniqID + "/event/txack",
			&as.TxAckEvent{},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{StrVal: `{}`},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId:   uniqID,
						Attr:     "ack_gateway_tx",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
					}, SkipToken: true,
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			bInpProto, err := proto.Marshal(lTest.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				bInpProto))

			// Don't stop the flow of execution (assert) to avoid leaving
			// messages orphaned in the queue.
			for _, res := range lTest.resVIn {
				select {
				case msg := <-globalVInDevSub.C():
					msg.Ack()
					t.Logf("Dev msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					assert.Equal(t, globalVInDevPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					assert.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.Point.TraceId
					// Normalize timestamps.
					assert.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
						testTimeout)
					res.Point.Ts = vIn.Point.Ts

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Errorf("\nExpect: %+v\nActual: %+v", lTest.resVIn,
							vIn)
					}
				case <-time.After(testTimeout):
					t.Error("Message timed out")
				}
			}

			if lTest.resPIn != nil {
				select {
				case msg := <-globalDInDataSub.C():
					msg.Ack()
					t.Logf("Data msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, globalDInPubTopic, msg.Topic())

					pIn := &message.DecoderIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), pIn))
					t.Logf("pIn: %+v", pIn)

					// Normalize generated trace ID.
					lTest.resPIn.TraceId = pIn.TraceId
					// Normalize timestamps.
					require.WithinDuration(t, time.Now(), pIn.Ts.AsTime(),
						testTimeout)
					lTest.resPIn.Ts = pIn.Ts

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(lTest.resPIn, pIn) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.resPIn,
							pIn)
					}
				case <-time.After(testTimeout):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestDecodeGatewaysDevicesError(t *testing.T) {
	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpPayl  []byte
	}{
		// Bad topic.
		{
			"lora", nil,
		},
		// Bad payload.
		{
			"lora/gateway/" + uniqID + "/event/up", []byte("ing-aaa"),
		},
		{
			"lora/application/1/device/" + uniqID + "/event/up",
			[]byte("ing-aaa"),
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, globalMQTTQueue.Publish(lTest.inpTopic,
				lTest.inpPayl))

			select {
			case msg := <-globalVInGWSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case msg := <-globalVInDevSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
