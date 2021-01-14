// +build !unit

package test

import (
	"crypto/rand"
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/api/api"
	"github.com/thingspect/atlas/internal/api/config"
	"github.com/thingspect/atlas/pkg/crypto"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/queue"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
)

var (
	globalOrgDAO  *org.DAO
	globalUserDAO *user.DAO

	globalPass string
	// globalHash is stored globally for test performance under -race.
	globalHash []byte

	globalNoAuthGRPCConn  *grpc.ClientConn
	globalAuthGRPCConn    *grpc.ClientConn
	globalAuthOrgID       string
	secondaryAuthGRPCConn *grpc.ClientConn

	globalPubTopic string
	globalPubSub   queue.Subber
)

func TestMain(m *testing.M) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalf("TestMain rand.Read: %v", err)
	}

	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgURI = testConfig.PgURI
	cfg.PWTKey = key

	cfg.NSQPubAddr = testConfig.NSQPubAddr
	cfg.NSQPubTopic += "-test-" + random.String(10)
	globalPubTopic = cfg.NSQPubTopic
	log.Printf("TestMain cfg.NSQPubTopic: %v", cfg.NSQPubTopic)

	// Set up API.
	a, err := api.New(cfg)
	if err != nil {
		log.Fatalf("TestMain api.New: %v", err)
	}

	// Serve connections.
	go func() {
		a.Serve()
	}()

	// Set up database connection.
	pg, err := postgres.New(cfg.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalUserDAO = user.NewDAO(pg)

	globalPass = random.String(10)
	globalHash, err = crypto.HashPass(globalPass)
	if err != nil {
		log.Fatalf("TestMain crypto.HashPass: %v", err)
	}
	log.Printf("globalPass, globalHash: %v, %s", globalPass, globalHash)

	// Build unauthenticated gRPC connection.
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	globalNoAuthGRPCConn, err = grpc.Dial(api.GRPCHost+api.GRPCPort, opts...)
	if err != nil {
		log.Fatalf("TestMain grpc.Dial: %v", err)
	}

	// Build authenticated gRPC connections.
	globalAuthOrgID, globalAuthGRPCConn, err = authGRPCConn(api.GRPCHost +
		api.GRPCPort)
	if err != nil {
		log.Fatalf("TestMain globalAuthGRPCConn authGRPCConn: %v", err)
	}

	_, secondaryAuthGRPCConn, err = authGRPCConn(api.GRPCHost +
		api.GRPCPort)
	if err != nil {
		log.Fatalf("TestMain secondaryAuthGRPCConn authGRPCConn: %v", err)
	}

	// Set up NSQ subscription to verify published messages. Use a unique
	// channel for each test run. This prevents failed tests from interfering
	// with the next run, but does require eventual cleaning.
	subChannel := api.ServiceName + "-test-" + random.String(10)
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, nil, subChannel,
		queue.DefaultNSQRequeueDelay)
	if err != nil {
		log.Fatalf("TestMain queue.NewNSQ: %v", err)
	}

	globalPubSub, err = nsq.Subscribe(cfg.NSQPubTopic)
	if err != nil {
		log.Fatalf("TestMain nsq.Subscribe: %v", err)
	}
	log.Printf("TestMain connected as NSQ sub channel: %v", subChannel)

	os.Exit(m.Run())
}
