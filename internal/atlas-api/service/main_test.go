// +build !integration

package service

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/thingspect/atlas/internal/atlas-api/crypto"
	"github.com/thingspect/atlas/pkg/test/random"
)

const (
	testTimeout = 2 * time.Second
	badUUID     = "..."
)

var (
	globalPass string
	// globalHash is stored globally for test performance under -race.
	globalHash []byte
)

func TestMain(m *testing.M) {
	var err error

	globalPass = random.String(10)
	globalHash, err = crypto.HashPass(globalPass)
	if err != nil {
		log.Fatalf("TestMain crypto.HashPass: %v", err)
	}
	log.Printf("globalPass, globalHash: %v, %s", globalPass, globalHash)

	os.Exit(m.Run())
}
