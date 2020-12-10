// +build !unit

package test

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/accumulator/accumulator"
	"github.com/thingspect/atlas/internal/accumulator/config"
	"github.com/thingspect/atlas/pkg/dao/datapoint"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var (
	globalVOutSubTopic string
	globalVOutQueue    queue.Queuer

	globalDPDAO *datapoint.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgURI = testConfig.PgURI

	cfg.NSQLookupAddrs = testConfig.NSQLookupAddrs
	cfg.NSQSubTopic += "-test-" + random.String(10)
	globalVOutSubTopic = cfg.NSQSubTopic
	log.Printf("TestMain cfg.NSQSubTopic: %v", cfg.NSQSubTopic)
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "accumulator-test-" + random.String(10)
	cfg.NSQPubAddr = testConfig.NSQPubAddr

	// Set up NSQ queue to publish test payloads.
	var err error
	globalVOutQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, "",
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain globalVOutQueue queue.NewNSQ: %v", err)
	}

	// Publish a throwaway message before subscribe to allow for discovery by
	// nsqlookupd.
	if err = globalVOutQueue.Publish(cfg.NSQSubTopic,
		[]byte("acc-aaa")); err != nil {
		log.Fatalf("TestMain globalVOutQueue.Publish: %v", err)
	}

	// Set up Accumulator.
	acc, err := accumulator.New(cfg)
	if err != nil {
		log.Fatalf("TestMain accumulator.New: %v", err)
	}

	// Serve connections.
	go func() {
		acc.Serve(cfg.Concurrency)
	}()

	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalDPDAO = datapoint.NewDAO(pg)

	os.Exit(m.Run())
}
