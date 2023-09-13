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

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v", lTest), func(t *testing.T) {
			t.Parallel()

			logDef.Debug("Debug")
			logDef.Debugf("Debugf: %v", lTest)
			logDef.Info("Info")
			logDef.Infof("Infof: %v", lTest)
			logDef.Error("Error")
			logDef.Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestDefaultConsole(t *testing.T) {
	SetDefault(NewConsole("DEBUG"))

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v", lTest), func(t *testing.T) {
			t.Parallel()

			Debug("Debug")
			Debugf("Debugf: %v", lTest)
			Info("Info")
			Infof("Infof: %v", lTest)
			Error("Error")
			Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestDefaultJSON(t *testing.T) {
	SetDefault(NewJSON("DEBUG"))

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v", lTest), func(t *testing.T) {
			t.Parallel()

			Debug("Debug")
			Debugf("Debugf: %v", lTest)
			Info("Info")
			Infof("Infof: %v", lTest)
			Error("Error")
			Errorf("Errorf: %v", lTest)
			// Do not test Fatal* due to os.Exit.
		})
	}
}

func TestDefaultWithField(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v with string", lTest), func(t *testing.T) {
			t.Parallel()

			logField := WithField(strconv.Itoa(lTest), random.String(10))
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
