//go:build !integration

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

func TestZlogWithField(t *testing.T) {
	t.Parallel()

	logger := newZlogJSON().WithField(random.String(10), random.String(10))
	t.Logf("logger: %#v", logger)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v with string", lTest), func(t *testing.T) {
			t.Parallel()

			logger.Debug("Debug")
			logger.Debugf("Debugf: %v", lTest)
			logger.Info("Info")
			logger.Infof("Infof: %v", lTest)
			logger.Error("Error")
			logger.Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}
