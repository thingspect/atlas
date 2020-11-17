// Package config provides configuration defaults and environment keys for
// tests.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "TEST_"

// Config holds settings used by test implementations.
type Config struct {
	PostgresURI    string
	NSQPubAddr     string
	NSQLookupAddrs []string
	MQTTAddr       string
}

// New instantiates a test Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		PostgresURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@localhost/postgres"),
		NSQPubAddr: config.String(pref+"NSQ_PUB_ADDR", "localhost:4150"),
		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"localhost:4161"}),
		MQTTAddr: config.String(pref+"MQTT_ADDR", "tcp://localhost:1883"),
	}
}
