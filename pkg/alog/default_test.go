//go:build !integration

package alog

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/thingspect/atlas/pkg/test/random"
)

func TestDefault(t *testing.T) {
	logDef := Default()
	t.Logf("logDef: %#v", logDef)

	for i := range 5 {
		t.Run(fmt.Sprintf("Can log %v", i), func(t *testing.T) {
			t.Parallel()

			logDef.Debug("Debug")
			logDef.Debugf("Debugf: %v", i)
			logDef.Info("Info")
			logDef.Infof("Infof: %v", i)
			logDef.Error("Error")
			logDef.Errorf("Errorf: %v", i)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestDefaultConsole(t *testing.T) {
	SetDefault(NewConsole("DEBUG"))

	for i := range 5 {
		t.Run(fmt.Sprintf("Can log %v", i), func(t *testing.T) {
			t.Parallel()

			Debug("Debug")
			Debugf("Debugf: %v", i)
			Info("Info")
			Infof("Infof: %v", i)
			Error("Error")
			Errorf("Errorf: %v", i)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestDefaultJSON(t *testing.T) {
	SetDefault(NewJSON("DEBUG"))

	for i := range 5 {
		t.Run(fmt.Sprintf("Can log %v", i), func(t *testing.T) {
			t.Parallel()

			Debug("Debug")
			Debugf("Debugf: %v", i)
			Info("Info")
			Infof("Infof: %v", i)
			Error("Error")
			Errorf("Errorf: %v", i)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestDefaultWithField(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can log %v with string", i), func(t *testing.T) {
			t.Parallel()

			logField := WithField(strconv.Itoa(i), random.String(10))
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
