package main

import (
	"github.com/thingspect/atlas/internal/accumulator/accumulator"
	"github.com/thingspect/atlas/internal/accumulator/config"
	"github.com/thingspect/atlas/pkg/alog"
)

func main() {
	cfg := config.New()

	alog.SetGlobal(alog.NewJSON().WithLevel(cfg.LogLevel).WithStr("service",
		accumulator.ServiceName))

	// Build Accumulator.
	acc, err := accumulator.New(cfg)
	if err != nil {
		alog.Fatalf("main accumulator.New: %v", err)
	}

	// Serve connections.
	acc.Serve(cfg.Concurrency)
}
