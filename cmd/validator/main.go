package main

import (
	"github.com/thingspect/atlas/internal/validator/config"
	"github.com/thingspect/atlas/internal/validator/validator"
	"github.com/thingspect/atlas/pkg/alog"
)

func main() {
	cfg := config.New()

	alog.SetGlobal(alog.NewJSON().WithLevel(cfg.LogLevel).WithStr("service",
		validator.ServiceName))

	// Build Validator.
	val, err := validator.New(cfg)
	if err != nil {
		alog.Fatalf("main validator.New: %v", err)
	}

	// Serve connections.
	val.Serve(cfg.Concurrency)
}
