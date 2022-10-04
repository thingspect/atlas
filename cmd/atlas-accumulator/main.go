// Package main starts the Accumulator service.
package main

import (
	"github.com/thingspect/atlas/internal/atlas-accumulator/accumulator"
	"github.com/thingspect/atlas/internal/atlas-accumulator/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON().WithLevel(cfg.LogLevel).WithField("service",
		accumulator.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, accumulator.ServiceName)

	// Build Accumulator.
	acc, err := accumulator.New(cfg)
	if err != nil {
		alog.Fatalf("main accumulator.New: %v", err)
	}

	// Serve connections.
	acc.Serve(cfg.Concurrency)
}
