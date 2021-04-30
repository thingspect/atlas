// Package config provides configuration values and defaults for the Decoder
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "DECODER_"

// Config holds settings used by the Decoder service.
type Config struct {
	LogLevel   string
	StatsDAddr string

	PgURI string

	NSQLookupAddrs []string
	NSQSubTopic    string
	NSQSubChannel  string

	NSQPubAddr  string
	NSQPubTopic string
	Concurrency int
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:   config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr: config.String(pref+"STATSD_ADDR", ""),

		PgURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),

		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"127.0.0.1:4161"}),
		NSQSubTopic:   config.String(pref+"NSQ_SUB_TOPIC", "DecoderIn"),
		NSQSubChannel: config.String(pref+"NSQ_SUB_CHANNEL", "decoder"),

		NSQPubAddr:  config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQPubTopic: config.String(pref+"NSQ_PUB_TOPIC", "ValidatorIn"),
		Concurrency: config.Int(pref+"CONCURRENCY", 5),
	}
}
