// Package main starts the Decoder service.
package main

import (
	"github.com/thingspect/atlas/internal/atlas-decoder/config"
	"github.com/thingspect/atlas/internal/atlas-decoder/decoder"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON().WithLevel(cfg.LogLevel).WithField("service",
		decoder.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, decoder.ServiceName)

	// Build Decoder.
	dec, err := decoder.New(cfg)
	if err != nil {
		alog.Fatalf("main decoder.New: %v", err)
	}

	// Serve connections.
	dec.Serve(cfg.Concurrency)
}
