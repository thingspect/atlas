// Package queue contains functions to publish and subscribe to queues.
package queue

// Messager defines the methods provided by a Message. Messages are not
// guaranteed to be thread-safe, and should only be accessed by their methods.
type Messager interface {
	Topic() string
	Payload() []byte
	Ack()
}

// Subber defines the methods provided by a Subscription.
type Subber interface {
	C() <-chan Messager
	Unsubscribe() error
}

// Queuer defines the methods provided by a Queue.
type Queuer interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string) (Subber, error)
	Disconnect()
}
