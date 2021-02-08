// Package config provides configuration values and defaults for the LoRa
// Ingestor service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "LORA_INGEST_"

// Config holds settings used by the LoRa Ingestor service.
type Config struct {
	LogLevel   string
	StatsDAddr string

	MQTTAddr   string
	MQTTUser   string
	MQTTPass   string
	MQTTShared bool

	NSQPubAddr      string
	NSQPubGWTopic   string
	NSQPubDevTopic  string
	NSQPubDataTopic string
	Concurrency     int
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel:   config.String(pref+"LOG_LEVEL", "DEBUG"),
		StatsDAddr: config.String(pref+"STATSD_ADDR", ""),

		MQTTAddr:   config.String(pref+"MQTT_ADDR", "tcp://127.0.0.1:1883"),
		MQTTUser:   config.String(pref+"MQTT_USER", "lora-ingestor"),
		MQTTPass:   config.String(pref+"MQTT_PASS", "notasecurepassword"),
		MQTTShared: config.Bool(pref+"MQTT_SHARED", true),

		NSQPubAddr:      config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQPubGWTopic:   config.String(pref+"NSQ_PUB_GW_TOPIC", "ValidatorIn"),
		NSQPubDevTopic:  config.String(pref+"NSQ_PUB_DEV_TOPIC", "ValidatorIn"),
		NSQPubDataTopic: config.String(pref+"NSQ_PUB_DATA_TOPIC", "ParserIn"),
		Concurrency:     config.Int(pref+"CONCURRENCY", 5),
	}
}
