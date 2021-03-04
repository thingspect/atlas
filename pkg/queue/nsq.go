package queue

import (
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	DefaultNSQRequeueDelay = 15 * time.Second
	nsqDisconnectTimeout   = 5 * time.Second
)

// nsqQueue contains methods to publish and subscribe to NSQ and implements the
// Queuer interface.
type nsqQueue struct {
	producer     *nsq.Producer
	pubAddr      string
	lookupAddrs  []string
	subChannel   string
	requeueDelay time.Duration
}

// Verify nsqQueue implements Queuer.
var _ Queuer = &nsqQueue{}

// NewNSQ builds a new Queuer and returns it and an error value. If lookupAddrs
// is nil, pubAddr is used for subscriptions. subChannel may be empty for a
// publish-only Queue. requeueDelay should usually be set to
// DefaultNSQRequeueDelay.
func NewNSQ(pubAddr string, lookupAddrs []string, subChannel string,
	requeueDelay time.Duration) (Queuer,
	error) {
	config := nsq.NewConfig()

	producer, err := nsq.NewProducer(pubAddr, config)
	if err != nil {
		return nil, err
	}

	if err = producer.Ping(); err != nil {
		return nil, err
	}

	return &nsqQueue{
		producer:     producer,
		pubAddr:      pubAddr,
		lookupAddrs:  lookupAddrs,
		subChannel:   subChannel,
		requeueDelay: requeueDelay,
	}, nil
}

// Publish publishes a message to a Queue and returns an error value.
func (n *nsqQueue) Publish(topic string, payload []byte) error {
	return n.producer.Publish(topic, payload)
}

// nsqSub contains methods to read from a subscription and implements the Subber
// interface.
type nsqSub struct {
	consumer *nsq.Consumer
	msgChan  chan Messager
}

// Verify nsqSub implements Subber.
var _ Subber = &nsqSub{}

// C returns the channel that carries a Subber's messages.
func (ns *nsqSub) C() <-chan Messager {
	return ns.msgChan
}

// Unsubscribe unsubscribes from a topic and returns an error value.
func (ns *nsqSub) Unsubscribe() error {
	ns.consumer.Stop()

	select {
	case <-ns.consumer.StopChan:
		close(ns.msgChan)

		return nil
	case <-time.After(nsqDisconnectTimeout):
		return ErrTimeout
	}
}

// nsqMessage contains methods to read from a message and implements the
// Messager interface.
type nsqMessage struct {
	topic string
	msg   *nsq.Message
}

// Verify nsqMessage implements Messager.
var _ Messager = &nsqMessage{}

// Topic returns the Messager's topic.
func (nm *nsqMessage) Topic() string {
	return nm.topic
}

// Payload returns the Messager's payload.
func (nm *nsqMessage) Payload() []byte {
	return nm.msg.Body
}

// Ack acknowledges successful processing of a Messager.
func (nm *nsqMessage) Ack() {
	nm.msg.Finish()
}

// Requeue requeues a Messager using a per-message backoff based on the
// number of attempts and DefaultRequeueDelay. This backoff is multiplicative
// and non-throttling. Requeue should only be used with transient failures that
// are likely to resolve.
func (nm *nsqMessage) Requeue() {
	nm.msg.RequeueWithoutBackoff(-1)
}

// Subscribe subscribes to a topic and returns a Subber and an error value.
func (n *nsqQueue) Subscribe(topic string) (Subber, error) {
	msgs := make(chan Messager)

	config := nsq.NewConfig()
	config.MaxInFlight = 10
	config.DefaultRequeueDelay = n.requeueDelay

	consumer, err := nsq.NewConsumer(topic, n.subChannel, config)
	if err != nil {
		return nil, err
	}

	consumer.SetLoggerLevel(nsq.LogLevelWarning)
	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		m.DisableAutoResponse()
		msgs <- &nsqMessage{topic: topic, msg: m}

		return nil
	}))

	if len(n.lookupAddrs) > 0 {
		err = consumer.ConnectToNSQLookupds(n.lookupAddrs)
	} else {
		err = consumer.ConnectToNSQD(n.pubAddr)
	}

	return &nsqSub{
		consumer: consumer,
		msgChan:  msgs,
	}, err
}

// Disconnect ends the connection to a Queue.
func (n *nsqQueue) Disconnect() {
	n.producer.Stop()
}
