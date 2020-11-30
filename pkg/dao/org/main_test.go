// +build !unit

package org

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/test/config"
)

var globalDAO *DAO

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connections.
	pg, err := postgres.New(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalDAO = NewDAO(pg)

	os.Exit(m.Run())
}
