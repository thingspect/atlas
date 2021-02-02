// +build !unit

package user

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/crypto"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/postgres"
	"github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 8 * time.Second

var (
	globalOrgDAO  *org.DAO
	globalUserDAO *DAO
	// globalHash is stored globally for test performance under -race.
	globalHash []byte
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
	globalUserDAO = NewDAO(pg)

	globalHash, err = crypto.HashPass(random.String(10))
	if err != nil {
		log.Fatalf("TestMain crypto.HashPass: %v", err)
	}
	log.Printf("globalHash: %s", globalHash)

	os.Exit(m.Run())
}
