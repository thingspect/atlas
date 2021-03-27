package main

import (
	"github.com/thingspect/atlas/internal/alerter/alerter"
	"github.com/thingspect/atlas/internal/alerter/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
)

func main() {
	cfg := config.New()

	alog.SetDefault(alog.NewJSON().WithLevel(cfg.LogLevel).WithStr("service",
		alerter.ServiceName))
	metric.SetStatsD(cfg.StatsDAddr, alerter.ServiceName)

	// Build Alerter.
	ale, err := alerter.New(cfg)
	if err != nil {
		alog.Fatalf("main alerter.New: %v", err)
	}

	// Serve connections.
	ale.Serve(cfg.Concurrency)
}
