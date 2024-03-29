//go:build !unit

package test

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/atlas-alerter/alerter"
	"github.com/thingspect/atlas/internal/atlas-alerter/config"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/alarm"
	"github.com/thingspect/atlas/pkg/dao/alert"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/rule"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

var (
	globalEOutSubTopic string
	globalAleQueue     queue.Queuer

	globalOrgDAO   *org.DAO
	globalRuleDAO  *rule.DAO
	globalUserDAO  *user.DAO
	globalAlarmDAO *alarm.DAO
	globalAleDAO   *alert.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgRwURI = testConfig.PgURI
	cfg.PgRoURI = testConfig.PgURI
	cfg.RedisHost = testConfig.RedisHost

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQLookupAddrs = testConfig.NSQLookupAddrs
	cfg.NSQSubTopic += "-test-" + random.String(10)
	globalEOutSubTopic = cfg.NSQSubTopic
	log.Printf("TestMain cfg.NSQSubTopic: %v", cfg.NSQSubTopic)
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "alerter-test-" + random.String(10)

	// Set up NSQ queue to publish test payloads.
	var err error
	globalAleQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, "")
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

	// Set up Alerter.
	ale, err := alerter.New(cfg)
	if err != nil {
		log.Fatalf("TestMain alerter.New: %v", err)
	}

	// Serve connections.
	go func() {
		ale.Serve(cfg.Concurrency)
	}()

	// Set up database connection.
	pg, err := dao.NewPgDB(cfg.PgRwURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg, pg)
	globalRuleDAO = rule.NewDAO(pg, pg)
	globalUserDAO = user.NewDAO(pg, pg)
	globalAlarmDAO = alarm.NewDAO(pg, pg)
	globalAleDAO = alert.NewDAO(pg, pg)

	os.Exit(m.Run())
}
