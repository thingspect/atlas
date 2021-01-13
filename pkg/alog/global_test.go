// +build !integration

package alog

import (
	"fmt"
	"testing"

	"github.com/thingspect/atlas/pkg/test/random"
)

func TestGlobal(t *testing.T) {
	logger := Global()
	t.Logf("logger: %#v", logger)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v", lTest), func(t *testing.T) {
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

func TestGlobalConsole(t *testing.T) {
	SetGlobal(NewConsole())

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

func TestGlobalJSON(t *testing.T) {
	SetGlobal(NewJSON())

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

func TestGlobalWithStr(t *testing.T) {
	t.Parallel()

	logger := WithStr(random.String(10), random.String(10))
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

func TestGlobalWithFields(t *testing.T) {
	t.Parallel()

	fields := map[string]interface{}{
		random.String(10): random.String(10),
		random.String(10): random.Intn(99),
	}
	logger := WithFields(fields)
	t.Logf("logger: %#v", logger)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can log %v with fields", lTest), func(t *testing.T) {
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
