//go:build !unit

package test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/atlas-lora-ingestor/config"
	"github.com/thingspect/atlas/internal/atlas-lora-ingestor/ingestor"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var (
	globalMQTTQueue queue.Queuer

	globalVInGWPubTopic  string
	globalVInDevPubTopic string
	globalDInPubTopic    string
	globalVInGWSub       queue.Subber
	globalVInDevSub      queue.Subber
	globalDInDataSub     queue.Subber
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.MQTTAddr = testConfig.MQTTAddr

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQGWPubTopic += "-test-" + random.String(10)
	globalVInGWPubTopic = cfg.NSQGWPubTopic
	log.Printf("TestMain cfg.NSQPubGWTopic: %v", cfg.NSQGWPubTopic)
	cfg.NSQDevPubTopic += "-test-" + random.String(10)
	globalVInDevPubTopic = cfg.NSQDevPubTopic
	log.Printf("TestMain cfg.NSQPubDevTopic: %v", cfg.NSQDevPubTopic)
	cfg.NSQDataPubTopic += "-test-" + random.String(10)
	globalDInPubTopic = cfg.NSQDataPubTopic
	log.Printf("TestMain cfg.NSQPubDataTopic: %v", cfg.NSQDataPubTopic)

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
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, subChannel)
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

	globalVInGWSub, err = nsq.Subscribe(cfg.NSQGWPubTopic)
	if err != nil {
		log.Fatalf("TestMain GW nsq.Subscribe: %v", err)
	}
	globalVInDevSub, err = nsq.Subscribe(cfg.NSQDevPubTopic)
	if err != nil {
		log.Fatalf("TestMain Dev nsq.Subscribe: %v", err)
	}
	globalDInDataSub, err = nsq.Subscribe(cfg.NSQDataPubTopic)
	if err != nil {
		log.Fatalf("TestMain Data nsq.Subscribe: %v", err)
	}
	log.Printf("TestMain connected as NSQ sub channel: %v", subChannel)

	os.Exit(m.Run())
}
