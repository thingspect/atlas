// +build !unit

package queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewMQTT(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inpAddr    string
		inpTimeout time.Duration
		err        error
	}{
		// Success.
		{testConfig.MQTTAddr, DefaultMQTTConnectTimeout, nil},
		// Wrong port.
		{"tcp://localhost:1884", 100 * time.Millisecond, ErrTimeout},
		// Unknown host.
		{"host-" + random.String(10), 100 * time.Millisecond, ErrTimeout},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can connect %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := NewMQTT(lTest.inpAddr, "", "",
				"testNewMQTT-"+random.String(10), lTest.inpTimeout)
			t.Logf("res, err: %+v, %v", res, err)
			if lTest.err == nil {
				require.NotNil(t, res)
			}
			require.Equal(t, lTest.err, err)
		})
	}
}

func TestMQTTPublish(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "",
		"testMQTTPublish-"+random.String(10), DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	require.NoError(t, mqtt.Publish("testMQTTPublish-"+random.String(10),
		[]byte(random.String(10))))
}

func TestMQTTSubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testMQTTSubscribe-" + random.String(10)
	payload := []byte(random.String(10))

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "",
		"testMQTTSubscribe-"+random.String(10), DefaultMQTTConnectTimeout)
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
	case <-time.After(250 * time.Millisecond):
		t.Error("Message timed out")
	}
}

func TestMQTTUnsubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testMQTTUnsubscribe-" + random.String(10)

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "",
		"testMQTTUnsubscribe-"+random.String(10), DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	sub, err := mqtt.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, sub.Unsubscribe())

	// Publish after unsubscribe to verify closed channel.
	require.NoError(t, mqtt.Publish("testMQTTUnsubscribe-"+random.String(10),
		[]byte(random.String(10))))

	select {
	case msg, ok := <-sub.C():
		t.Logf("msg, ok: %#v, %v", msg, ok)
		require.Nil(t, msg)
		require.False(t, ok)
	case <-time.After(250 * time.Millisecond):
		t.Error("Message timed out")
	}
}

func TestMQTTDisconnect(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "",
		"testMQTTDisconnect-"+random.String(10), DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	mqtt.Disconnect()
}
