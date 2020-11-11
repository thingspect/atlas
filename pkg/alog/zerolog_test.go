// +build !integration

package alog

import (
	"fmt"
	"testing"

	"github.com/thingspect/atlas/pkg/test/random"
)

func TestZlogWithLevel(t *testing.T) {
	t.Parallel()

	log := newZlogConsole()
	t.Logf("log: %#v", log)
	log.Debug("Debug (default)")

	tests := []string{
		"DEBUG",
		"info",
		"ERROR",
		"fatal",
		"BADDEBUG",
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can log %v", lTest), func(t *testing.T) {
			t.Parallel()

			logLevel := log.WithLevel(lTest)
			t.Logf("logLevel: %#v", logLevel)
			logLevel.Debug("Debug")
			logLevel.Debugf("Debugf: %#v", logLevel)
			logLevel.Info("Info")
			logLevel.Infof("Infof: %#v", logLevel)
			logLevel.Error("Error")
			logLevel.Errorf("Errorf: %#v", logLevel)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestZlogAllFields(t *testing.T) {
	t.Parallel()

	log := newZlogJSON().WithStr(random.String(10), random.String(10))
	t.Logf("log: %#v", log)
	log.Debug("WithStr")

	fields := map[string]interface{}{
		random.String(10): random.String(10),
		random.String(10): random.Intn(99),
	}
	log = log.WithFields(fields)
	t.Logf("log: %#v", log)
	log.Debug("WithStr + WithFields")
}
