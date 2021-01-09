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

func TestNewNSQ(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inpPubAddr     string
		inpLookupAddrs []string
		err            string
	}{
		// Success.
		{testConfig.NSQPubAddr, nil, ""},
		{testConfig.NSQPubAddr, testConfig.NSQLookupAddrs, ""},
		// Wrong port.
		{"127.0.0.1:4152", nil, "connect: connection refused"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can connect %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := NewNSQ(lTest.inpPubAddr, lTest.inpLookupAddrs,
				"testNewNSQ-"+random.String(10), DefaultNSQRequeueDelay)
			t.Logf("res, err: %+v, %v", res, err)
			if lTest.err == "" {
				require.NotNil(t, res)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}

func TestNSQPublish(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	nsq, err := NewNSQ(testConfig.NSQPubAddr, testConfig.NSQLookupAddrs,
		"testNSQPublish-"+random.String(10), DefaultNSQRequeueDelay)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	require.NoError(t, nsq.Publish("testNSQPublish-"+random.String(10),
		[]byte(random.String(10))))

	nsq, err = NewNSQ(testConfig.NSQPubAddr, nil, "", DefaultNSQRequeueDelay)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	require.NoError(t, nsq.Publish("testNSQPublish-"+random.String(10),
		[]byte(random.String(10))))
}

func TestNSQSubscribeLookup(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testNSQSubscribeLookup-" + random.String(10)
	payload := []byte(random.String(10))

	nsq, err := NewNSQ(testConfig.NSQPubAddr, testConfig.NSQLookupAddrs,
		"testNSQSubscribeLookup-"+random.String(10), DefaultNSQRequeueDelay)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	// Publish before subscribe to allow for discovery by nsqlookupd.
	require.NoError(t, nsq.Publish(topic, payload))

	sub, err := nsq.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, payload, msg.Payload())
	case <-time.After(5 * time.Second):
		t.Fatal("Message timed out")
	}
}

func TestNSQSubscribePub(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testNSQSubscribePub-" + random.String(10)
	payload := []byte(random.String(10))

	nsq, err := NewNSQ(testConfig.NSQPubAddr, nil,
		"testNSQSubscribePub-"+random.String(10), DefaultNSQRequeueDelay)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	sub, err := nsq.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, nsq.Publish(topic, payload))

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, payload, msg.Payload())
	case <-time.After(5 * time.Second):
		t.Fatal("Message timed out")
	}
}

func TestNSQUnsubscribe(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testNSQUnsubscribe-" + random.String(10)

	nsq, err := NewNSQ(testConfig.NSQPubAddr, testConfig.NSQLookupAddrs,
		"testNSQUnsubscribe-"+random.String(10), DefaultNSQRequeueDelay)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	sub, err := nsq.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	require.NoError(t, sub.Unsubscribe())

	// Publish after unsubscribe to verify closed channel.
	require.NoError(t, nsq.Publish("testNSQUnsubscribe-"+random.String(10),
		[]byte(random.String(10))))

	select {
	case msg, ok := <-sub.C():
		t.Logf("msg, ok: %#v, %v", msg, ok)
		require.Nil(t, msg)
		require.False(t, ok)
	case <-time.After(5 * time.Second):
		t.Fatal("Message timed out")
	}
}

func TestNSQRequeue(t *testing.T) {
	t.Parallel()

	testConfig := config.New()
	topic := "testNSQRequeue-" + random.String(10)
	payload := []byte(random.String(10))

	nsq, err := NewNSQ(testConfig.NSQPubAddr, testConfig.NSQLookupAddrs,
		"testNSQRequeue-"+random.String(10), time.Millisecond)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	// Publish before subscribe to allow for discovery by nsqlookupd.
	require.NoError(t, nsq.Publish(topic, payload))

	sub, err := nsq.Subscribe(topic)
	t.Logf("sub, err: %+v, %v", sub, err)
	require.NoError(t, err)

	select {
	case msg := <-sub.C():
		msg.Requeue()
		t.Logf("Requeue msg.Topic, msg.Payload: %v, %x", msg.Topic(),
			msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, payload, msg.Payload())
	case <-time.After(5 * time.Second):
		t.Fatal("Requeue message timed out")
	}

	select {
	case msg := <-sub.C():
		msg.Ack()
		t.Logf("Ack msg.Topic, msg.Payload: %v, %x", msg.Topic(), msg.Payload())
		require.Equal(t, topic, msg.Topic())
		require.Equal(t, payload, msg.Payload())
	case <-time.After(10 * time.Second):
		t.Fatal("Ack message timed out")
	}
}

func TestNSQDisconnect(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	nsq, err := NewNSQ(testConfig.NSQPubAddr, testConfig.NSQLookupAddrs,
		"testNSQDisconnect-"+random.String(10), DefaultNSQRequeueDelay)
	t.Logf("nsq, err: %+v, %v", nsq, err)
	require.NoError(t, err)

	nsq.Disconnect()
}
