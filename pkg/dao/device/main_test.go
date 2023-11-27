//go:build !unit

package device

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/config"
)

var (
	globalOrgDAO      *org.DAO
	globalDevDAO      *DAO
	globalDevDAOCache *DAO
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
	globalDevDAO = NewDAO(pg, pg, nil, 0)
	globalDevDAOCache = NewDAO(pg, pg, cache.NewMemory(), time.Minute)

	os.Exit(m.Run())
}
