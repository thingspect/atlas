// +build !unit

package test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/mqtt-ingestor/config"
	"github.com/thingspect/atlas/internal/mqtt-ingestor/ingestor"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var globalMQTT queue.Queuer
var globalParser queue.Subber

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.MQTTAddr = testConfig.MQTTAddr
	cfg.NSQPubAddr = testConfig.NSQPubAddr

	// Set up Ingestor.
	ing, err := ingestor.New(cfg)
	if err != nil {
		log.Fatalf("TestMain ingestor.New: %v", err)
	}

	// Serve connections.
	go func() {
		ing.Serve(cfg.ParserConcurrency)
	}()

	// Set up MQTT client connection to publish test payloads.
	clientID := fmt.Sprintf("%s-test-%s", ingestor.ServiceName,
		random.String(10))
	globalMQTT, err = queue.NewMQTT(cfg.MQTTAddr, cfg.MQTTUser, cfg.MQTTPass,
		clientID, queue.DefaultMQTTConnectTimeout)
	if err != nil {
		log.Fatalf("TestMain queue.NewMQTT: %v", err)
	}
	log.Printf("TestMain connected as MQTT client: %v", clientID)

	// Set up queue to verify published messages. Use a unique channel for each
	// test run. This prevents failed tests from interfering with the next run,
	// but does require eventual cleaning.
	subChannel := "mqtt-ingestor-test-" + random.String(10)
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, subChannel,
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

	globalParser, err = nsq.Subscribe("ValidatorIn")
	if err != nil {
		log.Fatalf("TestMain nsq.Subscribe: %v", err)
	}

	os.Exit(m.Run())
}
