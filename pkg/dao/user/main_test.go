// +build !unit

package user

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/config"
)

const testTimeout = 8 * time.Second

var (
	globalOrgDAO  *org.DAO
	globalUserDAO *DAO
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
	globalUserDAO = NewDAO(pg)

	os.Exit(m.Run())
}
