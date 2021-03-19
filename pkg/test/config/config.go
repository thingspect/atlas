// Package config provides configuration values and defaults for tests. This
// package must not be used outside of tests.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "TEST_"

// Config holds settings used by test implementations.
type Config struct {
	PgURI     string
	RedisHost string

	NSQPubAddr     string
	NSQLookupAddrs []string

	MQTTAddr string
	MQTTUser string
	MQTTPass string
}

// New instantiates a test Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		PgURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		RedisHost: config.String(pref+"REDIS_HOST", "127.0.0.1"),

		NSQPubAddr: config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"127.0.0.1:4161"}),

		MQTTAddr: config.String(pref+"MQTT_ADDR", "tcp://127.0.0.1:1883"),
		MQTTUser: config.String(pref+"MQTT_USER", "atlas_test"),
		MQTTPass: config.String(pref+"MQTT_PASS", "notasecurepassword"),
	}
}
