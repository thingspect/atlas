// +build !unit

package device

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/test/config"
)

var globalDAO *DAO
var globalOrgDAO *org.DAO

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connections.
	pg, err := postgres.New(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalDAO = NewDAO(pg)
	globalOrgDAO = org.NewDAO(pg)

	os.Exit(m.Run())
}
