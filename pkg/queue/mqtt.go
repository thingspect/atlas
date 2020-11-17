package queue

import (
	"errors"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	DefaultMQTTConnectTimeout = 5 * time.Second
	mqttPublishTimeout        = 30 * time.Second
)

var ErrTimeout = errors.New("queue: timed out")

// mqttQueue contains methods to publish and subscribe to MQTT and implements
// the Queuer interface.
type mqttQueue struct {
	client         mqtt.Client
	connectTimeout time.Duration
}

// Verify mqttQueue implements Queuer.
var _ Queuer = &mqttQueue{}

// NewMQTT builds a new Queue and returns a reference to it and an error value.
// connectTimeout should usually be set to DefaultMQTTConnectTimeout.
func NewMQTT(addr, user, pass, clientID string,
	connectTimeout time.Duration) (Queuer, error) {
	// Build client options and assign to a client.
	opts := mqtt.NewClientOptions().
		AddBroker(addr).
		SetUsername(user).
		SetPassword(pass).
		SetClientID(clientID).
		SetMaxReconnectInterval(connectTimeout)
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if ok := token.WaitTimeout(connectTimeout); !ok {
		return nil, ErrTimeout
	}

	return &mqttQueue{
		client:         client,
		connectTimeout: connectTimeout,
	}, token.Error()
}

// Publish publishes a message to a Queue and returns an error value.
func (m *mqttQueue) Publish(topic string, payload []byte) error {
	token := m.client.Publish(topic, 1, false, payload)
	if ok := token.WaitTimeout(mqttPublishTimeout); !ok {
		return ErrTimeout
	}

	return token.Error()
}

// mqttSub contains methods to read from a subscription and implements the
// Subber interface.
type mqttSub struct {
	mqtt    *mqttQueue
	topic   string
	msgChan chan Messager
}

// Verify mqttSub implements Subber.
var _ Subber = &mqttSub{}

// C returns the channel that carries a Subber's messages.
func (ms *mqttSub) C() <-chan Messager {
	return ms.msgChan
}

// Unsubscribe unsubscribes to a topic and returns an error value.
func (ms *mqttSub) Unsubscribe() error {
	token := ms.mqtt.client.Unsubscribe(ms.topic)
	if ok := token.WaitTimeout(ms.mqtt.connectTimeout); !ok {
		return ErrTimeout
	}

	if err := token.Error(); err != nil {
		return err
	}

	close(ms.msgChan)

	return nil
}

// mqttMessage contains methods to read from a message and implements the
// Messager interface.
type mqttMessage struct {
	mqtt.Message
}

// Verify mqttMessage implements Messager.
var _ Messager = &mqttMessage{}

// Requeue is not supported.
func (mm *mqttMessage) Requeue() {
	log.Fatal("Requeue unsupported")
}

// Subscribe subscribes to a topic and returns a Subber and an error value.
func (m *mqttQueue) Subscribe(topic string) (Subber, error) {
	msgs := make(chan Messager)

	token := m.client.Subscribe(topic, 1,
		func(client mqtt.Client, msg mqtt.Message) {
			msgs <- &mqttMessage{Message: msg}
		})
	if ok := token.WaitTimeout(m.connectTimeout); !ok {
		return nil, ErrTimeout
	}

	return &mqttSub{
		mqtt:    m,
		topic:   topic,
		msgChan: msgs,
	}, token.Error()
}

// Disconnect ends the connection to a Queue.
func (m *mqttQueue) Disconnect() {
	m.client.Disconnect(uint(m.connectTimeout / time.Millisecond))
}
