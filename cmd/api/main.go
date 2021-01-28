package main

import (
	"github.com/thingspect/atlas/internal/api/api"
	"github.com/thingspect/atlas/internal/api/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetGlobal(alog.NewJSON().WithLevel(cfg.LogLevel).WithStr("service",
		api.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, api.ServiceName)

	// Build API.
	a, err := api.New(cfg)
	if err != nil {
		alog.Fatalf("main api.New: %v", err)
	}

	// Serve connections.
	a.Serve()
}
