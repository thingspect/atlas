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

func TestZlogWithStr(t *testing.T) {
	t.Parallel()

	logEntry := newZlogJSON().WithStr(random.String(10), random.String(10))
	t.Logf("logEntry: %#v", logEntry)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v with string", lTest), func(t *testing.T) {
			t.Parallel()

			logEntry.Debug("Debug")
			logEntry.Debugf("Debugf: %v", lTest)
			logEntry.Info("Info")
			logEntry.Infof("Infof: %v", lTest)
			logEntry.Error("Error")
			logEntry.Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestZlogWithFields(t *testing.T) {
	t.Parallel()

	fields := map[string]interface{}{
		random.String(10): random.String(10),
		random.String(10): random.Intn(99),
	}
	logEntry := newZlogJSON().WithFields(fields)
	t.Logf("logEntry: %#v", logEntry)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v with fields", lTest), func(t *testing.T) {
			t.Parallel()

			logEntry.Debug("Debug")
			logEntry.Debugf("Debugf: %v", lTest)
			logEntry.Info("Info")
			logEntry.Infof("Infof: %v", lTest)
			logEntry.Error("Error")
			logEntry.Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}
