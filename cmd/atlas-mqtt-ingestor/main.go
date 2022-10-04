// Package main starts the MQTT Ingestor service.
package main

import (
	"github.com/thingspect/atlas/internal/atlas-mqtt-ingestor/config"
	"github.com/thingspect/atlas/internal/atlas-mqtt-ingestor/ingestor"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON().WithLevel(cfg.LogLevel).WithField("service",
		ingestor.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, ingestor.ServiceName)

	// Build Ingestor.
	ing, err := ingestor.New(cfg)
	if err != nil {
		alog.Fatalf("main ingestor.New: %v", err)
	}

	// Serve connections.
	ing.Serve(cfg.Concurrency)
}
