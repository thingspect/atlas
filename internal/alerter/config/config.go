// Package config provides configuration values and defaults for the Alerter
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "ALERTER_"

// Config holds settings used by the Alerter service.
type Config struct {
	LogLevel   string
	StatsDAddr string

	PgURI     string
	RedisHost string

	NSQLookupAddrs []string
	NSQSubTopic    string
	NSQSubChannel  string

	NSQPubAddr  string
	Concurrency int

	AppAPIKey     string
	SMSAccountSID string
	SMSAuthToken  string
	SMSPhone      string
	EmailAPIKey   string
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:   config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr: config.String(pref+"STATSD_ADDR", ""),

		PgURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		RedisHost: config.String(pref+"REDIS_HOST", "127.0.0.1"),

		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"127.0.0.1:4161"}),
		NSQSubTopic:   config.String(pref+"NSQ_SUB_TOPIC", "EventerOut"),
		NSQSubChannel: config.String(pref+"NSQ_SUB_CHANNEL", "alerter"),

		NSQPubAddr:  config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		Concurrency: config.Int(pref+"CONCURRENCY", 5),

		AppAPIKey: config.String(pref+"APP_API_KEY", ""),
		SMSAccountSID: config.String(pref+"SMS_ACCOUNT_SID",
			"ACdefd5896d2acc15061cae169bb9d4836"),
		SMSAuthToken: config.String(pref+"SMS_AUTH_TOKEN", ""),
		SMSPhone:     config.String(pref+"SMS_PHONE", "+15125432462"),
		EmailAPIKey:  config.String(pref+"EMAIL_API_KEY", ""),
	}
}
