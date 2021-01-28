// Package config provides configuration values and defaults for the API
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "API_"

// Config holds settings used by the API service.
type Config struct {
	LogLevel   string
	StatsDAddr string

	PgURI  string
	PWTKey []byte

	NSQPubAddr  string
	NSQPubTopic string
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:   config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr: config.String(pref+"STATSD_ADDR", ""),

		PgURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		PWTKey: config.ByteSlice(pref + "PWT_KEY"),

		NSQPubAddr:  config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQPubTopic: config.String(pref+"NSQ_PUB_TOPIC", "ValidatorIn"),
	}
}
