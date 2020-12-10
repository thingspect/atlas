// Package config provides configuration values and defaults for the Accumulator
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "ACCUMULATOR_"

// Config holds settings used by the Accumulator service.
type Config struct {
	LogLevel string

	PgURI string

	NSQLookupAddrs []string
	NSQSubTopic    string
	NSQSubChannel  string

	NSQPubAddr  string
	Concurrency int
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel: config.String(pref+"LOG_LEVEL", "DEBUG"),

		PgURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas"),

		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"127.0.0.1:4161"}),
		NSQSubTopic:   config.String(pref+"NSQ_SUB_TOPIC", "ValidatorOut"),
		NSQSubChannel: config.String(pref+"NSQ_SUB_CHANNEL", "accumulator"),

		NSQPubAddr:  config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		Concurrency: config.Int(pref+"CONCURRENCY", 5),
	}
}
