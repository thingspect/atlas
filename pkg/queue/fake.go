package queue

// fakeQueue contains methods to publish and subscribe to a fake Queue and
// implements the Queuer interface.
type fakeQueue struct {
	msgChan chan Messager
}

// Verify fakeQueue implements Queuer.
var _ Queuer = &fakeQueue{}

// NewFake builds a new Queuer using a single buffered channel and returns it.
// The channel will never be closed, even on unsubscribe.
func NewFake() Queuer {
	return &fakeQueue{
		msgChan: make(chan Messager, 10),
	}
}

// fakeMessage contains methods to read from a message and implements the
// Messager interface.
type fakeMessage struct {
	topic   string
	payload []byte
}

// Verify fakeMessage implements Messager.
var _ Messager = &fakeMessage{}

// Topic returns the Messager's topic.
func (fm *fakeMessage) Topic() string {
	return fm.topic
}

// Payload returns the Messager's payload.
func (fm *fakeMessage) Payload() []byte {
	return fm.payload
}

// Ack mocks acknowledging successful processing of a Messager.
func (fm *fakeMessage) Ack() {}

// Requeue mocks requeueing a Messager.
func (fm *fakeMessage) Requeue() {}

// Publish publishes a message to a Queue and returns an error value.
func (f *fakeQueue) Publish(topic string, payload []byte) error {
	f.msgChan <- &fakeMessage{topic: topic, payload: payload}

	return nil
}

// Prime primes a Queue topic by publishing a single-byte message, with value
// Prime, for the purpose of being discarded.
func (f *fakeQueue) Prime(topic string) error {
	return f.Publish(topic, []byte{Prime})
}

// fakeSub contains methods to read from a subscription and implements the
// Subber interface.
type fakeSub struct {
	msgChan chan Messager
}

// Verify fakeSub implements Subber.
var _ Subber = &fakeSub{}

// C returns the channel that carries a Subber's messages.
func (fs *fakeSub) C() <-chan Messager {
	return fs.msgChan
}

// Unsubscribe mocks unsubscribing from a topic and returns an error value.
func (fs *fakeSub) Unsubscribe() error {
	return nil
}

// Subscribe subscribes to the default topic and returns a Subber and an error
// value. The provided topic is discarded.
func (f *fakeQueue) Subscribe(_ string) (Subber, error) {
	return &fakeSub{msgChan: f.msgChan}, nil
}

// Disconnect mocks ending the connection to a Queue.
func (f *fakeQueue) Disconnect() {}
