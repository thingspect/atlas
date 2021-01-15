// +build !unit

package datapoint

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/test/config"
)

var (
	globalDPDAO  *DAO
	globalOrgDAO *org.DAO
	globalDevDAO *device.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := postgres.New(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalDPDAO = NewDAO(pg)
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = device.NewDAO(pg)

	os.Exit(m.Run())
}
