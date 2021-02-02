// +build !unit

package org

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/test/config"
)

const testTimeout = 8 * time.Second

var globalOrgDAO *DAO

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := postgres.New(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalOrgDAO = NewDAO(pg)

	os.Exit(m.Run())
}
