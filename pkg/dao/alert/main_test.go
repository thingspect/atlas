// +build !unit

package alert

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/config"
)

var (
	globalOrgDAO   *org.DAO
	globalDevDAO   *device.DAO
	globalAlertDAO *DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := dao.NewPgDB(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = device.NewDAO(pg)
	globalAlertDAO = NewDAO(pg)

	os.Exit(m.Run())
}
