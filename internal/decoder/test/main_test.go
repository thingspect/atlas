// +build !unit

package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/internal/decoder/config"
	"github.com/thingspect/atlas/internal/decoder/decoder"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 6 * time.Second

var (
	globalDInSubTopic string
	globalDInQueue    queue.Queuer

	globalDecoderPubTopic string
	globalDecoderSub      queue.Subber

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
	globalDInSubTopic = cfg.NSQSubTopic
	log.Printf("TestMain cfg.NSQSubTopic: %v", cfg.NSQSubTopic)
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "decoder-test-" + random.String(10)

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQPubTopic += "-test-" + random.String(10)
	globalDecoderPubTopic = cfg.NSQPubTopic
	log.Printf("TestMain cfg.NSQPubTopic: %v", cfg.NSQPubTopic)

	// Set up NSQ queue to publish test payloads.
	var err error
	globalDInQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, "",
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain globalDInQueue queue.NewNSQ: %v", err)
	}

	// Publish a throwaway message before subscribe to allow for discovery by
	// nsqlookupd.
	if err = globalDInQueue.Publish(cfg.NSQSubTopic,
		[]byte("dec-aaa")); err != nil {
		log.Fatalf("TestMain globalDInQueue.Publish: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	log.Print("TestMain published throwaway message")

	// Set up Decoder.
	dec, err := decoder.New(cfg)
	if err != nil {
		log.Fatalf("TestMain decoder.New: %v", err)
	}

	// Serve connections.
	go func() {
		dec.Serve(cfg.Concurrency)
	}()

	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = device.NewDAO(pg)

	// Set up NSQ subscription to verify published messages.
	decoderQueue, err := queue.NewNSQ(cfg.NSQPubAddr, nil, cfg.NSQSubChannel,
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain decoderQueue queue.NewNSQ: %v", err)
	}

	globalDecoderSub, err = decoderQueue.Subscribe(cfg.NSQPubTopic)
	if err != nil {
		log.Fatalf("TestMain decoderQueue.Subscribe: %v", err)
	}
	log.Printf("TestMain connected as NSQ sub channel: %v", cfg.NSQSubChannel)

	os.Exit(m.Run())
}
