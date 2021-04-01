// +build !integration

package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestFakePublish(t *testing.T) {
	t.Parallel()

	fake := NewFake()
	t.Logf("fake: %+v", fake)

	require.NoError(t, fake.Publish("testFakePublish-"+random.String(10),
		random.Bytes(10)))
}

func TestFakeSubscribe(t *testing.T) {
	t.Parallel()

	topic := "testFakeSubscribe-" + random.String(10)
	payload := random.Bytes(10)

	fake := NewFake()
	t.Logf("fake: %+v", fake)

	sub, err := fake.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, fake.Publish(topic, payload))

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, payload, msg.Payload())
	case <-time.After(2 * time.Second):
		t.Fatal("Message timed out")
	}
}

func TestFakePrime(t *testing.T) {
	t.Parallel()

	topic := "testFakePrime-" + random.String(10)

	fake := NewFake()
	t.Logf("fake: %+v", fake)

	sub, err := fake.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, fake.Prime(topic))

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, []byte{Prime}, msg.Payload())
	case <-time.After(2 * time.Second):
		t.Fatal("Message timed out")
	}
}

func TestFakeUnsubscribe(t *testing.T) {
	t.Parallel()

	topic := "testFakeSubscribePub-" + random.String(10)

	fake := NewFake()
	t.Logf("fake: %+v", fake)

	sub, err := fake.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, sub.Unsubscribe())
}

func TestFakeDisconnect(t *testing.T) {
	t.Parallel()

	fake := NewFake()
	t.Logf("fake: %+v", fake)

	fake.Disconnect()
}
