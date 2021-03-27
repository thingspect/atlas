// +build !unit

package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/internal/alerter/alerter"
	"github.com/thingspect/atlas/internal/alerter/config"
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
	cfg.PgURI = testConfig.PgURI
	cfg.RedisHost = testConfig.RedisHost

	cfg.NSQLookupAddrs = testConfig.NSQLookupAddrs
	cfg.NSQSubTopic += "-test-" + random.String(10)
	globalEOutSubTopic = cfg.NSQSubTopic
	log.Printf("TestMain cfg.NSQSubTopic: %v", cfg.NSQSubTopic)
	// Use a unique channel for each test run. This prevents failed tests from
	// interfering with the next run, but does require eventual cleaning.
	cfg.NSQSubChannel = "alerter-test-" + random.String(10)
	cfg.NSQPubAddr = testConfig.NSQPubAddr

	// Set up NSQ queue to publish test payloads.
	var err error
	globalAleQueue, err = queue.NewNSQ(cfg.NSQPubAddr, nil, "",
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

	// Publish a throwaway message before subscribe to allow for discovery by
	// nsqlookupd.
	if err = globalAleQueue.Publish(cfg.NSQSubTopic,
		[]byte("ale-aaa")); err != nil {
		log.Fatalf("TestMain globalAleQueue.Publish: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	log.Print("TestMain published throwaway message")

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
	pg, err := dao.NewPgDB(cfg.PgURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalRuleDAO = rule.NewDAO(pg)
	globalUserDAO = user.NewDAO(pg)
	globalAlarmDAO = alarm.NewDAO(pg)
	globalAleDAO = alert.NewDAO(pg)

	os.Exit(m.Run())
}
