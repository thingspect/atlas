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
		{random.String(10), 100 * time.Millisecond, ErrTimeout},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can connect %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := NewMQTT(lTest.inpAddr, "", "", random.String(10),
				lTest.inpTimeout)
			t.Logf("res, err: %+v, %v", res, err)
			if lTest.err == nil {
				require.NotNil(t, res)
			}
			require.Equal(t, lTest.err, err)
		})
	}
}

func TestPublish(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "", random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	require.NoError(t, mqtt.Publish(random.String(10),
		[]byte(random.String(10))))
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := random.String(10)
	payload := []byte(random.String(10))

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "", random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	sub, err := mqtt.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, mqtt.Publish(topic, payload))

	msg := <-sub.C()
	t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
	require.Equal(t, topic, msg.Topic())
	require.Equal(t, payload, msg.Payload())
}

func TestUnubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := random.String(10)

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "", random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	sub, err := mqtt.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	unsubTopic := sub.Topic()
	require.Equal(t, topic, unsubTopic)
	require.NoError(t, mqtt.Unsubscribe(unsubTopic))
}

func TestDisconnect(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	mqtt, err := NewMQTT(testConfig.MQTTAddr, "", "", random.String(10),
		DefaultMQTTConnectTimeout)
	t.Logf("mqtt, err: %+v, %v", mqtt, err)
	require.NoError(t, err)

	mqtt.Disconnect()
}
