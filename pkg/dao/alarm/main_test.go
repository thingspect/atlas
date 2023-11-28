//go:build !unit

package alarm

import (
	"log"
	"os"
	"testing"

	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/rule"
	"github.com/thingspect/atlas/pkg/test/config"
)

var (
	globalOrgDAO   *org.DAO
	globalRuleDAO  *rule.DAO
	globalAlarmDAO *DAO
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
	globalRuleDAO = rule.NewDAO(pg, pg)
	globalAlarmDAO = NewDAO(pg, pg)

	os.Exit(m.Run())
}
