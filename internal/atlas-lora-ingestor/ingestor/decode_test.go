//go:build !integration

package ingestor

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/protobuf/proto"
)

func TestDecodeGateways(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		res      []*message.ValidatorIn
	}{
		{
			"lora/us915_0/gateway/" + uniqID + "/event/up", &gw.UplinkFrame{
				RxInfo: &gw.UplinkRxInfo{Rssi: -74},
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
			"lora/us915_0/gateway/" + uniqID + "/event/stats", &gw.GatewayStats{
				RxPacketsReceivedOk: 2,
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_gateway",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"rxPacketsReceivedOk":2}`,
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
			"lora/us915_0/gateway/" + uniqID + "/event/ack", &gw.DownlinkTxAck{
				Items: []*gw.DownlinkTxAckItem{{Status: gw.TxAckStatus_OK}},
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
			"lora/us915_0/gateway/" + uniqID + "/event/exec",
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
			"lora/us915_0/gateway/" + uniqID + "/state/conn", &gw.ConnState{
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
		t.Run(fmt.Sprintf("Can decode %+v", test), func(t *testing.T) {
			t.Parallel()

			mqttGWQueue := queue.NewFake()
			mqttGWSub, err := mqttGWQueue.Subscribe("")
			require.NoError(t, err)

			ingQueue := queue.NewFake()
			vInSub, err := ingQueue.Subscribe("")
			require.NoError(t, err)
			vInGWPubTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttGWSub: mqttGWSub,

				ingQueue:      ingQueue,
				vInGWPubTopic: vInGWPubTopic,
			}
			go func() {
				ing.decodeGateways()
			}()

			bInpProto, err := proto.Marshal(test.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, mqttGWQueue.Publish(test.inpTopic, bInpProto))

			for _, res := range test.res {
				select {
				case msg := <-vInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, vInGWPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.GetPoint().GetTraceId()
					// Normalize timestamp.
					require.WithinDuration(t, time.Now(),
						vIn.GetPoint().GetTs().AsTime(), 2*time.Second)
					res.Point.Ts = vIn.GetPoint().GetTs()

					require.EqualExportedValues(t, res, vIn)
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
		{"lora/us915_0/gateway/" + uniqID + "/event/up", []byte("ing-aaa")},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can decode %+v", test), func(t *testing.T) {
			t.Parallel()

			mqttGWQueue := queue.NewFake()
			mqttGWSub, err := mqttGWQueue.Subscribe("")
			require.NoError(t, err)

			ingQueue := queue.NewFake()
			vInSub, err := ingQueue.Subscribe("")
			require.NoError(t, err)
			vInGWPubTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttGWSub: mqttGWSub,

				ingQueue:      ingQueue,
				vInGWPubTopic: vInGWPubTopic,
			}
			go func() {
				ing.decodeGateways()
			}()

			require.NoError(t, mqttGWQueue.Publish(test.inpTopic,
				test.inpPayl))

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

func TestDecodeDevices(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)
	devAddr := random.String(16)

	bData := random.Bytes(10)
	b64Data := base64.StdEncoding.EncodeToString(bData)
	t.Logf("b64Data: %v", b64Data)

	tests := []struct {
		inpTopic string
		inpProto proto.Message
		res      []*message.ValidatorIn
	}{
		{
			"lora/application/1/device/" + uniqID + "/event/up",
			&integration.UplinkEvent{Data: bData},
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
		},
		{
			"lora/application/2/device/" + uniqID + "/event/join",
			&integration.JoinEvent{DevAddr: devAddr},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: fmt.Sprintf(`{"devAddr":"%s"}`, devAddr),
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
						UniqId: uniqID, Attr: "devaddr",
						ValOneof: &common.DataPoint_StrVal{StrVal: devAddr},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/application/3/device/" + uniqID + "/event/ack",
			&integration.AckEvent{Acknowledged: true},
			[]*message.ValidatorIn{
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
			},
		},
		{
			"lora/application/4/device/" + uniqID + "/event/log",
			&integration.LogEvent{Code: integration.LogCode_OTAA},
			[]*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "raw_device",
						ValOneof: &common.DataPoint_StrVal{
							StrVal: `{"code":"OTAA"}`,
						},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "log_level",
						ValOneof: &common.DataPoint_StrVal{StrVal: "INFO"},
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "log_code",
						ValOneof: &common.DataPoint_StrVal{StrVal: "OTAA"},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/application/5/device/" + uniqID + "/event/txack",
			&integration.TxAckEvent{},
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
						Attr:     "tx_queued",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
					}, SkipToken: true,
				},
			},
		},
		{
			"lora/application/6/device/" + uniqID + "/event/status",
			&integration.StatusEvent{},
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
						Attr:     "ext_power",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
					}, SkipToken: true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can decode %+v", test), func(t *testing.T) {
			t.Parallel()

			mqttDevQueue := queue.NewFake()
			mqttDevSub, err := mqttDevQueue.Subscribe("")
			require.NoError(t, err)

			ingQueue := queue.NewFake()
			vInSub, err := ingQueue.Subscribe("")
			require.NoError(t, err)
			vInDevPubTopic := "topic-" + random.String(10)
			dInPubTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttDevSub: mqttDevSub,

				ingQueue:       ingQueue,
				vInDevPubTopic: vInDevPubTopic,
				dInPubTopic:    dInPubTopic,
			}
			go func() {
				ing.decodeDevices()
			}()

			bInpProto, err := proto.Marshal(test.inpProto)
			require.NoError(t, err)
			t.Logf("bInpProto: %s", bInpProto)

			require.NoError(t, mqttDevQueue.Publish(test.inpTopic, bInpProto))

			for _, res := range test.res {
				select {
				case msg := <-vInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, vInDevPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Normalize generated trace ID.
					res.Point.TraceId = vIn.GetPoint().GetTraceId()
					// Normalize timestamp.
					require.WithinDuration(t, time.Now(),
						vIn.GetPoint().GetTs().AsTime(), 2*time.Second)
					res.Point.Ts = vIn.GetPoint().GetTs()

					require.EqualExportedValues(t, res, vIn)
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
		{
			"lora", nil,
		},
		// Bad payload.
		{
			"lora/application/1/device/" + uniqID + "/event/up",
			[]byte("ing-aaa"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can decode %+v", test), func(t *testing.T) {
			t.Parallel()

			mqttDevQueue := queue.NewFake()
			mqttDevSub, err := mqttDevQueue.Subscribe("")
			require.NoError(t, err)

			ingQueue := queue.NewFake()
			vInSub, err := ingQueue.Subscribe("")
			require.NoError(t, err)
			vInDevPubTopic := "topic-" + random.String(10)
			dInPubTopic := "topic-" + random.String(10)

			ing := Ingestor{
				mqttDevSub: mqttDevSub,

				ingQueue:       ingQueue,
				vInDevPubTopic: vInDevPubTopic,
				dInPubTopic:    dInPubTopic,
			}
			go func() {
				ing.decodeDevices()
			}()

			require.NoError(t, mqttDevQueue.Publish(test.inpTopic,
				test.inpPayl))

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
