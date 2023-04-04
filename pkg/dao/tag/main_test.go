//go:build !unit

package tag

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/test/config"
)

var (
	globalOrgDAO  *org.DAO
	globalDevDAO  *device.DAO
	globalUserDAO *user.DAO
	globalTagDAO  *DAO
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
	globalDevDAO = device.NewDAO(pg, nil, 0)
	globalUserDAO = user.NewDAO(pg)
	globalTagDAO = NewDAO(pg)

	os.Exit(m.Run())
}
