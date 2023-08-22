// Package main starts the Eventer service.
package main

import (
	"github.com/thingspect/atlas/internal/atlas-eventer/config"
	"github.com/thingspect/atlas/internal/atlas-eventer/eventer"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON(cfg.LogLevel).WithField("service",
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
