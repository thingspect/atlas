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
		t.Run(fmt.Sprintf("Can log %v", test), func(t *testing.T) {
			t.Parallel()

			logLevel := newStlogConsole(test)
			t.Logf("logLevel: %#v", logLevel)

			logLevel.Debug("Debug")
			logLevel.Debugf("Debugf: %v and above", test)
			logLevel.Info("Info")
			logLevel.Infof("Infof: %v and above", test)
			logLevel.Error("Error")
			logLevel.Errorf("Errorf: %v and above", test)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestStlogWithField(t *testing.T) {
	t.Parallel()

	logger := newStlogJSON("DEBUG")
	t.Logf("logger: %#v", logger)

	for i := range 5 {
		t.Run(fmt.Sprintf("Can log %v with string", i), func(t *testing.T) {
			t.Parallel()

			logField := logger.WithField(strconv.Itoa(i), random.String(10))
			t.Logf("logField: %#v", logField)

			logField.Debug("Debug")
			logField.Debugf("Debugf: %v", i)
			logField.Info("Info")
			logField.Infof("Infof: %v", i)
			logField.Error("Error")
			logField.Errorf("Errorf: %v", i)
			// Do not test Fatal* due to os.Exit.
		})
	}
}
