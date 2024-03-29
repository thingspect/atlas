// Package ingestor provides functions used to run the MQTT Ingestor service.
package ingestor

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/thingspect/atlas/internal/atlas-mqtt-ingestor/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/queue"
)

// ServiceName provides consistent naming, including logs and metrics.
const ServiceName = "mqtt-ingestor"

// Constants used for MQTT topic configuration.
const (
	sharedPref  = "$share/v1group/"
	mqttV1Topic = "v1/#"
)

// Ingestor holds references to the message broker connections.
type Ingestor struct {
	mqttSub queue.Subber

	ingQueue    queue.Queuer
	vInPubTopic string
}

// New builds a new Ingestor and returns a reference to it and an error value.
func New(cfg *config.Config) (*Ingestor, error) {
	// Build the MQTT connection for consuming.
	id, err := rand.Int(rand.Reader, big.NewInt(99999))
	if err != nil {
		return nil, err
	}
	clientID := fmt.Sprintf("%s-%d", ServiceName, id)

	mqtt, err := queue.NewMQTT(cfg.MQTTAddr, cfg.MQTTUser, cfg.MQTTPass,
		clientID, queue.DefaultMQTTConnectTimeout)
	if err != nil {
		return nil, err
	}

	// Subscribe to the topic.
	topic := mqttV1Topic
	if cfg.MQTTShared {
		topic = sharedPref + topic
	}
	mqttSub, err := mqtt.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	// Build the NSQ connection for publishing.
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, "")
	if err != nil {
		return nil, err
	}

	return &Ingestor{
		mqttSub: mqttSub,

		ingQueue:    nsq,
		vInPubTopic: cfg.NSQPubTopic,
	}, nil
}

// Serve starts the message decoders.
func (ing *Ingestor) Serve(concurrency int) {
	for range concurrency {
		go ing.decodeMessages()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := ing.mqttSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve ing.mqttSub.Unsubscribe: %v", err)
	}
	ing.ingQueue.Disconnect()
}
