// +build !unit

package datapoint

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/config"
)

const testTimeout = 4 * time.Second

var (
	globalDPDAO  *DAO
	globalOrgDAO *org.DAO
	globalDevDAO *device.DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := dao.NewPgDB(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalDPDAO = NewDAO(pg)
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = device.NewDAO(pg)

	os.Exit(m.Run())
}
