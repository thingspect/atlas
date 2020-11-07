// Package config provides configuration defaults and environment keys for
// tests.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "TEST_"

// Config holds settings used by test implementations.
type Config struct {
	PostgresURI string
}

// New instantiates a test Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		PostgresURI: config.String(pref+"PGURI",
			"postgres://postgres:postgres@localhost/postgres"),
	}
}
