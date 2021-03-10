package main

import (
	"github.com/thingspect/atlas/internal/eventer/config"
	"github.com/thingspect/atlas/internal/eventer/eventer"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON().WithLevel(cfg.LogLevel).WithStr("service",
		eventer.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, eventer.ServiceName)

	// Build Eventer.
	ev, err := eventer.New(cfg)
	if err != nil {
		alog.Fatalf("main eventer.New: %v", err)
	}

	// Serve connections.
	ev.Serve(cfg.Concurrency)
}
