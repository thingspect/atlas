//go:build !unit

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
	globalOrgDAO *org.DAO
	globalDevDAO *device.DAO
	globalAleDAO *DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := dao.NewPgDB(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain dao.NewPgDB: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg, pg)
	globalDevDAO = device.NewDAO(pg, pg, nil, 0)
	globalAleDAO = NewDAO(pg, pg)

	os.Exit(m.Run())
}
