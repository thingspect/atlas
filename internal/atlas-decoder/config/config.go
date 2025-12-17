// Package config provides configuration values and defaults for the Decoder
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "DECODER_"

// Config holds settings used by the Decoder service.
type Config struct {
	LogLevel    string
	StatsDAddr  string
	Concurrency int

	PgRwURI string
	PgRoURI string

	NSQPubAddr     string
	NSQLookupAddrs []string
	NSQSubTopic    string
	NSQSubChannel  string
	NSQPubTopic    string
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:    config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr:  config.String(pref+"STATSD_ADDR", ""),
		Concurrency: config.Int(pref+"CONCURRENCY", 5),

		PgRwURI: config.String(pref+"PG_RW_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		PgRoURI: config.String(pref+"PG_RO_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),

		NSQPubAddr: config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"127.0.0.1:4161"}),
		NSQSubTopic:   config.String(pref+"NSQ_SUB_TOPIC", "DecoderIn"),
		NSQSubChannel: config.String(pref+"NSQ_SUB_CHANNEL", "decoder"),
		NSQPubTopic:   config.String(pref+"NSQ_PUB_TOPIC", "ValidatorIn"),
	}
}
