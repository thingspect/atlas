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

var (
	globalOrgDAO *org.DAO
	globalDevDAO *DAO
)

func TestMain(m *testing.M) {
	// Set up Config.
	testConfig := config.New()

	// Set up database connection.
	pg, err := postgres.New(testConfig.PgURI)
	if err != nil {
		log.Fatalf("TestMain postgres.New: %v", err)
	}
	globalOrgDAO = org.NewDAO(pg)
	globalDevDAO = NewDAO(pg)

	os.Exit(m.Run())
}
