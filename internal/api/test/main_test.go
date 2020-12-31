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

	globalGRPCConn *grpc.ClientConn
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

	// Build gRPC connection.
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	globalGRPCConn, err = grpc.Dial(api.GRPCHost+api.GRPCPort, opts...)
	if err != nil {
		log.Fatalf("TestMain grpc.Dial: %v", err)
	}

	os.Exit(m.Run())
}
