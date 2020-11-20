package main

import (
	"github.com/thingspect/atlas/internal/mqtt-ingestor/config"
	"github.com/thingspect/atlas/internal/mqtt-ingestor/ingestor"
	"github.com/thingspect/atlas/pkg/alog"
)

func main() {
	cfg := config.New()

	alog.SetGlobal(alog.NewJSON().WithLevel(cfg.LogLevel).WithStr("service",
		ingestor.ServiceName))

	// Build Ingestor.
	ingestor, err := ingestor.New(cfg)
	if err != nil {
		alog.Fatalf("main ingestor.New: %v", err)
	}

	// Serve connections.
	ingestor.Serve(cfg.ParserConcurrency)
}
