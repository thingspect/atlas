//go:build !integration

package alog

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/thingspect/atlas/pkg/test/random"
)

func TestStlogLevel(t *testing.T) {
	t.Parallel()

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

			logLevel := newStlogConsole(lTest)
			t.Logf("logLevel: %#v", logLevel)

			logLevel.Debug("Debug")
			logLevel.Debugf("Debugf: %v and above", lTest)
			logLevel.Info("Info")
			logLevel.Infof("Infof: %v and above", lTest)
			logLevel.Error("Error")
			logLevel.Errorf("Errorf: %v and above", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestStlogWithField(t *testing.T) {
	t.Parallel()

	logger := newStlogJSON("DEBUG")
	t.Logf("logger: %#v", logger)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v with string", lTest), func(t *testing.T) {
			t.Parallel()

			logField := logger.WithField(strconv.Itoa(lTest), random.String(10))
			t.Logf("logField: %#v", logField)

			logField.Debug("Debug")
			logField.Debugf("Debugf: %v", lTest)
			logField.Info("Info")
			logField.Infof("Infof: %v", lTest)
			logField.Error("Error")
			logField.Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}
