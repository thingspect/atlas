// Package main starts the Validator service.
package main

import (
	"github.com/thingspect/atlas/internal/atlas-validator/config"
	"github.com/thingspect/atlas/internal/atlas-validator/validator"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON().WithLevel(cfg.LogLevel).WithField("service",
		validator.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, validator.ServiceName)

	// Build Validator.
	val, err := validator.New(cfg)
	if err != nil {
		alog.Fatalf("main validator.New: %v", err)
	}

	// Serve connections.
	val.Serve(cfg.Concurrency)
}
