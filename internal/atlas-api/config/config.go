// Package config provides configuration values and defaults for the API
// service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "API_"

// Config holds settings used by the API service.
type Config struct {
	LogLevel   string
	StatsDAddr string

	PgRwURI   string
	PgRoURI   string
	RedisHost string

	NSQPubAddr  string
	NSQPubTopic string

	AppAPIKey    string
	SMSKeyID     string
	SMSKeySecret string

	LoRaAddr      string
	LoRaAPIKey    string
	LoRaTenantID  string
	LoRaAppID     string
	LoRaDevProfID string

	PWTKey  []byte
	APIHost string
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:   config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr: config.String(pref+"STATSD_ADDR", ""),

		PgRwURI: config.String(pref+"PG_RW_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		PgRoURI: config.String(pref+"PG_RO_URI",
			"postgres://postgres:postgres@127.0.0.1/atlas_test"),
		RedisHost: config.String(pref+"REDIS_HOST", "127.0.0.1"),

		NSQPubAddr:  config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQPubTopic: config.String(pref+"NSQ_PUB_TOPIC", "ValidatorIn"),

		AppAPIKey: config.String(pref+"APP_API_KEY", ""),
		SMSKeyID: config.String(pref+"SMS_KEY_ID",
			"SKb62d2a1320d85ad96b07a90fe92e051e"),
		SMSKeySecret: config.String(pref+"SMS_KEY_SECRET", ""),

		LoRaAddr:      config.String(pref+"LORA_ADDR", ""),
		LoRaAPIKey:    config.String(pref+"LORA_API_KEY", ""),
		LoRaTenantID:  config.String(pref+"LORA_TENANT_ID", ""),
		LoRaAppID:     config.String(pref+"LORA_APP_ID", ""),
		LoRaDevProfID: config.String(pref+"LORA_DEV_PROF_ID", ""),

		PWTKey:  config.ByteSlice(pref + "PWT_KEY"),
		APIHost: config.String(pref+"API_HOST", ""),
	}
}
