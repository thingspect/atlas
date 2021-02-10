// +build !unit

package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/internal/lora-ingestor/config"
	"github.com/thingspect/atlas/internal/lora-ingestor/ingestor"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 6 * time.Second

var (
	globalMQTTQueue queue.Queuer

	globalDecoderPubGWTopic   string
	globalDecoderPubDevTopic  string
	globalDecoderPubDataTopic string
	globalDecoderGWSub        queue.Subber
	globalDecoderDevSub       queue.Subber
	globalDecoderDataSub      queue.Subber
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.MQTTAddr = testConfig.MQTTAddr

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQPubGWTopic += "-test-" + random.String(10)
	globalDecoderPubGWTopic = cfg.NSQPubGWTopic
	log.Printf("TestMain cfg.NSQPubGWTopic: %v", cfg.NSQPubGWTopic)
	cfg.NSQPubDevTopic += "-test-" + random.String(10)
	globalDecoderPubDevTopic = cfg.NSQPubDevTopic
	log.Printf("TestMain cfg.NSQPubDevTopic: %v", cfg.NSQPubDevTopic)
	cfg.NSQPubDataTopic += "-test-" + random.String(10)
	globalDecoderPubDataTopic = cfg.NSQPubDataTopic
	log.Printf("TestMain cfg.NSQPubDataTopic: %v", cfg.NSQPubDataTopic)

	// Set up MQTT client connection to publish test payloads.
	var err error
	clientID := fmt.Sprintf("%s-test-%s", ingestor.ServiceName,
		random.String(10))
	globalMQTTQueue, err = queue.NewMQTT(cfg.MQTTAddr, cfg.MQTTUser,
		cfg.MQTTPass, clientID, queue.DefaultMQTTConnectTimeout)
	if err != nil {
		log.Fatalf("TestMain queue.NewMQTT: %v", err)
	}
	log.Printf("TestMain connected as MQTT client: %v", clientID)

	// Set up Ingestor.
	ing, err := ingestor.New(cfg)
	if err != nil {
		log.Fatalf("TestMain ingestor.New: %v", err)
	}

	// Serve connections.
	go func() {
		ing.Serve(cfg.Concurrency)
	}()

	// Set up NSQ subscription to verify published messages. Use a unique
	// channel for each test run. This prevents failed tests from interfering
	// with the next run, but does require eventual cleaning.
	subChannel := ingestor.ServiceName + "-test-" + random.String(10)
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, subChannel,
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

	globalDecoderGWSub, err = nsq.Subscribe(cfg.NSQPubGWTopic)
	if err != nil {
		log.Fatalf("TestMain GW nsq.Subscribe: %v", err)
	}
	globalDecoderDevSub, err = nsq.Subscribe(cfg.NSQPubDevTopic)
	if err != nil {
		log.Fatalf("TestMain Dev nsq.Subscribe: %v", err)
	}
	globalDecoderDataSub, err = nsq.Subscribe(cfg.NSQPubDataTopic)
	if err != nil {
		log.Fatalf("TestMain Data nsq.Subscribe: %v", err)
	}
	log.Printf("TestMain connected as NSQ sub channel: %v", subChannel)

	os.Exit(m.Run())
}
