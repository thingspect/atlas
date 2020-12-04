// +build !unit

package test

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/validator/config"
	"github.com/thingspect/atlas/internal/validator/validator"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var (
	globalVInSubTopic string
	globalVInQueue    queue.Queuer

	globalVOutPubTopic string
	globalVOutSub      queue.Subber

	globalOrgDAO *org.DAO
	globalDevDAO *device.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgURI = testConfig.PgURI

	cfg.NSQLookupAddrs = testConfig.NSQLookupAddrs
	cfg.NSQSubTopic += "-test-" + random.String(10)
	globalVInSubTopic = cfg.NSQSubTopic
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "validator-test-" + random.String(10)

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQPubTopic += "-test-" + random.String(10)
	globalVOutPubTopic = cfg.NSQPubTopic

	// Set up NSQ queue to publish test payloads.
	var err error
	globalVInQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, "",
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain globalVInQueue queue.NewNSQ: %v", err)
	}

	// Publish a throwaway message before subscribe to allow for discovery by
	// nsqlookupd.
	if err = globalVInQueue.Publish(cfg.NSQSubTopic,
		[]byte("val-aaa")); err != nil {
		log.Fatalf("TestMain globalVInQueue.Publish: %v", err)
	}

	// Set up Validator.
	val, err := validator.New(cfg)
	if err != nil {
		log.Fatalf("TestMain validator.New: %v", err)
	}

	// Serve connections.
	go func() {
		val.Serve(cfg.Concurrency)
	}()

	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = device.NewDAO(pg)

	// Set up NSQ subscription to verify published messages.
	vOutQueue, err := queue.NewNSQ(cfg.NSQPubAddr, nil, cfg.NSQSubChannel,
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain vOutQueue queue.NewNSQ: %v", err)
	}

	globalVOutSub, err = vOutQueue.Subscribe(cfg.NSQPubTopic)
	if err != nil {
		log.Fatalf("TestMain vOutQueue.Subscribe: %v", err)
	}
	log.Printf("TestMain connected as NSQ sub channel: %v", cfg.NSQSubChannel)

	os.Exit(m.Run())
}
