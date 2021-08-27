//go:build !unit

package test

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/atlas-accumulator/accumulator"
	"github.com/thingspect/atlas/internal/atlas-accumulator/config"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/datapoint"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var (
	globalVOutSubTopic string
	globalAccQueue     queue.Queuer

	globalOrgDAO *org.DAO
	globalDevDAO *device.DAO
	globalDPDAO  *datapoint.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgURI = testConfig.PgURI

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQLookupAddrs = testConfig.NSQLookupAddrs
	cfg.NSQSubTopic += "-test-" + random.String(10)
	globalVOutSubTopic = cfg.NSQSubTopic
	log.Printf("TestMain cfg.NSQSubTopic: %v", cfg.NSQSubTopic)
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "accumulator-test-" + random.String(10)

	// Set up NSQ queue to publish test payloads.
	var err error
	globalAccQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, "")
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
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
	pg, err := dao.NewPgDB(cfg.PgURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = device.NewDAO(pg)
	globalDPDAO = datapoint.NewDAO(pg)

	os.Exit(m.Run())
}
