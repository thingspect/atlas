// Package config provides configuration values and defaults for the Alerter
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "ALERTER_"

// Config holds settings used by the Alerter service.
type Config struct {
	LogLevel    string
	StatsDAddr  string
	Concurrency int

	PgURI     string
	RedisHost string

	NSQPubAddr     string
	NSQLookupAddrs []string
	NSQSubTopic    string
	NSQSubChannel  string

	AppAPIKey    string
	SMSKeyID     string
	SMSAccountID string
	SMSKeySecret string
	SMSPhone     string
	EmailDomain  string
	EmailAPIKey  string
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:    config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr:  config.String(pref+"STATSD_ADDR", ""),
		Concurrency: config.Int(pref+"CONCURRENCY", 5),

		PgURI: config.String(pref+"PG_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		RedisHost: config.String(pref+"REDIS_HOST", "127.0.0.1"),

		NSQPubAddr: config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQLookupAddrs: config.StringSlice(pref+"NSQ_LOOKUP_ADDRS",
			[]string{"127.0.0.1:4161"}),
		NSQSubTopic:   config.String(pref+"NSQ_SUB_TOPIC", "EventerOut"),
		NSQSubChannel: config.String(pref+"NSQ_SUB_CHANNEL", "alerter"),

		AppAPIKey: config.String(pref+"APP_API_KEY", ""),
		SMSKeyID: config.String(pref+"SMS_KEY_ID",
			"SKb62d2a1320d85ad96b07a90fe92e051e"),
		SMSAccountID: config.String(pref+"SMS_ACCOUNT_ID",
			"ACdefd5896d2acc15061cae169bb9d4836"),
		SMSKeySecret: config.String(pref+"SMS_KEY_SECRET", ""),
		SMSPhone:     config.String(pref+"SMS_PHONE", "+15125432462"),
		EmailDomain:  config.String(pref+"EMAIL_DOMAIN", "mg.thingspect.com"),
		EmailAPIKey:  config.String(pref+"EMAIL_API_KEY", ""),
	}
}
