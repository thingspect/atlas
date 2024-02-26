//go:build !unit

package queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 5 * time.Second

func TestNewMQTT(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inpAddr    string
		inpTimeout time.Duration
		err        string
	}{
		// Success.
		{
			testConfig.MQTTAddr, DefaultMQTTConnectTimeout, "",
		},
		// Wrong port.
		{
			"tcp://127.0.0.1:1884", DefaultMQTTConnectTimeout,
			"connect: connection refused",
		},
		// Unknown host.
		{
			"host-" + random.String(10) + ":1883", time.Millisecond,
			ErrTimeout.Error(),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can connect %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := NewMQTT(test.inpAddr, testConfig.MQTTUser,
				testConfig.MQTTPass, "testNewMQTT-"+random.String(10),
				test.inpTimeout)
			t.Logf("res, err: %+v, %#v", res, err)
			if test.err == "" {
				require.NotNil(t, res)
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}

func TestMQTTPublish(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	mqtt, err := NewMQTT(testConfig.MQTTAddr, testConfig.MQTTUser,
		testConfig.MQTTPass, "testMQTTPublish-"+random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	require.NoError(t, mqtt.Publish("testMQTTPublish-"+random.String(10),
		random.Bytes(10)))
}

func TestMQTTSubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testMQTTSubscribe-" + random.String(10)
	payload := random.Bytes(10)

	mqtt, err := NewMQTT(testConfig.MQTTAddr, testConfig.MQTTUser,
		testConfig.MQTTPass, "testMQTTSubscribe-"+random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	sub, err := mqtt.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, mqtt.Publish(topic, payload))

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, payload, msg.Payload())
	case <-time.After(testTimeout):
		t.Fatal("Message timed out")
	}
}

func TestMQTTPrime(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testMQTTPrime-" + random.String(10)

	mqtt, err := NewMQTT(testConfig.MQTTAddr, testConfig.MQTTUser,
		testConfig.MQTTPass, "testMQTTPrime-"+random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	sub, err := mqtt.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, mqtt.Prime(topic))

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, []byte{Prime}, msg.Payload())
	case <-time.After(testTimeout):
		t.Fatal("Message timed out")
	}
}

func TestMQTTUnsubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testMQTTUnsubscribe-" + random.String(10)

	mqtt, err := NewMQTT(testConfig.MQTTAddr, testConfig.MQTTUser,
		testConfig.MQTTPass, "testMQTTUnsubscribe-"+random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	sub, err := mqtt.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, sub.Unsubscribe())

	// Publish after unsubscribe to verify closed channel.
	require.NoError(t, mqtt.Publish("testMQTTUnsubscribe-"+random.String(10),
		random.Bytes(10)))

	select {
	case msg, ok := <-sub.C():
		t.Logf("msg, ok: %#v, %v", msg, ok)
		require.Nil(t, msg)
		require.False(t, ok)
	case <-time.After(testTimeout):
		t.Fatal("Message timed out")
	}
}

func TestMQTTDisconnect(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	mqtt, err := NewMQTT(testConfig.MQTTAddr, testConfig.MQTTUser,
		testConfig.MQTTPass, "testMQTTDisconnect-"+random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	mqtt.Disconnect()
}
