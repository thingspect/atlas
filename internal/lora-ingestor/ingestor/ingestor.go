// Package ingestor provides functions used to run the LoRa Ingestor service.
package ingestor

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/thingspect/atlas/internal/lora-ingestor/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/queue"
)

const (
	ServiceName  = "lora-ingestor"
	sharedPref   = "$share/loragroup/"
	mqttGWTopic  = "lora/gateway/+/event/+"
	mqttDevTopic = "lora/application/+/device/+/event/+"
)

// Ingestor holds references to the message broker connections.
type Ingestor struct {
	mqttGWSub  queue.Subber
	mqttDevSub queue.Subber

	ingQueue       queue.Queuer
	vInGWPubTopic  string
	vInDevPubTopic string
	dInPubTopic    string
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

	// Subscribe to the gateway topic.
	topic := mqttGWTopic
	if cfg.MQTTShared {
		topic = sharedPref + topic
	}
	mqttGWSub, err := mqtt.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	// Subscribe to the device topic.
	topic = mqttDevTopic
	if cfg.MQTTShared {
		topic = sharedPref + topic
	}
	mqttDevSub, err := mqtt.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	// Build the NSQ connection for publishing.
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, "",
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		return nil, err
	}

	return &Ingestor{
		mqttGWSub:  mqttGWSub,
		mqttDevSub: mqttDevSub,

		ingQueue:       nsq,
		vInGWPubTopic:  cfg.NSQGWPubTopic,
		vInDevPubTopic: cfg.NSQDevPubTopic,
		dInPubTopic:    cfg.NSQDataPubTopic,
	}, nil
}

// Serve starts the message decoders.
func (ing *Ingestor) Serve(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go ing.decodeGateways()
		go ing.decodeDevices()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := ing.mqttGWSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve ing.mqttGWSub.Unsubscribe: %v", err)
	}
	if err := ing.mqttDevSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve ing.mqttDevSub.Unsubscribe: %v", err)
	}
	ing.ingQueue.Disconnect()
}
