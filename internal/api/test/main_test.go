// +build !unit

package test

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/internal/api/api"
	"github.com/thingspect/atlas/internal/api/config"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/postgres"
	testconfig "github.com/thingspect/atlas/pkg/test/config"
	"google.golang.org/grpc"
)

var (
	globalOrgDAO   *org.DAO
	globalGRPCConn *grpc.ClientConn
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := testconfig.New()
	cfg := config.New()
	cfg.PgURI = testConfig.PgURI

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

	// Build gRPC connection.
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	globalGRPCConn, err = grpc.Dial("127.0.0.1:50051", opts...)
	if err != nil {
		log.Fatalf("TestMain grpc.Dial: %v", err)
	}

	os.Exit(m.Run())
}
