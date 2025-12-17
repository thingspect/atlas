//go:build !unit

package test

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/atlas-decoder/config"
	"github.com/thingspect/atlas/internal/atlas-decoder/decoder"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var (
	globalDInSubTopic string
	globalDecQueue    queue.Queuer

	globalVInPubTopic string
	globalVInSub      queue.Subber

	globalOrgDAO *org.DAO
	globalDevDAO *device.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgRwURI = testConfig.PgURI
	cfg.PgRoURI = testConfig.PgURI

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQLookupAddrs = testConfig.NSQLookupAddrs
	cfg.NSQSubTopic += "-test-" + random.String(10)
	globalDInSubTopic = cfg.NSQSubTopic
	log.Printf("TestMain cfg.NSQSubTopic: %v", cfg.NSQSubTopic)
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "decoder-test-" + random.String(10)

	cfg.NSQPubTopic += "-test-" + random.String(10)
	globalVInPubTopic = cfg.NSQPubTopic
	log.Printf("TestMain cfg.NSQPubTopic: %v", cfg.NSQPubTopic)

	// Set up NSQ queue to publish test payloads.
	var err error
	globalDecQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, cfg.NSQSubChannel)
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

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
	pg, err := dao.NewPgDB(cfg.PgRwURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg, pg)
	globalDevDAO = device.NewDAO(pg, pg, nil, 0)

	// Set up NSQ subscription to verify published messages.
	globalVInSub, err = globalDecQueue.Subscribe(cfg.NSQPubTopic)
	if err != nil {
		log.Fatalf("TestMain globalDecQueue.Subscribe: %v", err)
	}
	log.Printf("TestMain connected as NSQ sub channel: %v", cfg.NSQSubChannel)

	os.Exit(m.Run())
}
