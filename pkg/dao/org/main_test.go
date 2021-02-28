// +build !unit

package org

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/config"
)

const testTimeout = 8 * time.Second

var globalOrgDAO *DAO

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := dao.NewPgDB(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = NewDAO(pg)

	os.Exit(m.Run())
}
