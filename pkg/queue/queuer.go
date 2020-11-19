// Package queue provdes functions to publish and subscribe to queues.
package queue

// Messager defines the methods provided by a Message. Messages are not
// guaranteed to be thread-safe, and should only be accessed by their methods.
type Messager interface {
	// Topic returns the Messager's topic.
	Topic() string
	// Payload returns the Messager's payload.
	Payload() []byte
	// Ack acknowledges successful processing of a Messager.
	Ack()
	// Requeue requeues a Messager using a per-message backoff based on the
	// number of attempts. Requeue should only be used with transient failures
	// that are likely to resolve. Requeue is not supported by all Messager
	// types, specifically MQTT.
	Requeue()
}

// Subber defines the methods provided by a Subscription.
type Subber interface {
	// C returns the channel that carries a Subber's messages.
	C() <-chan Messager
	// Unsubscribe unsubscribes to a topic and returns an error value.
	Unsubscribe() error
}

// Queuer defines the methods provided by a Queue.
type Queuer interface {
	// Publish publishes a message to a Queue and returns an error value.
	Publish(topic string, payload []byte) error
	// Subscribe subscribes to a topic and returns a Subber and an error value.
	Subscribe(topic string) (Subber, error)
	// Disconnect ends the connection to a Queue.
	Disconnect()
}
