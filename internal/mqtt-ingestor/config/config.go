// Package config provides configuration values and defaults for the MQTT
// Ingestor service.
package config

import "github.com/thingspect/atlas/pkg/config"

const pref = "MQTT_INGEST_"

// Config holds settings used by the MQTT Ingestor service.
type Config struct {
	LogLevel string

	MQTTAddr   string
	MQTTUser   string
	MQTTPass   string
	MQTTShared bool

	NSQPubAddr  string
	NSQPubTopic string
	Concurrency int
}

// New instantiates a service Config, parses the environment, and returns it.
func New() *Config {
	return &Config{
		LogLevel: config.String(pref+"LOG_LEVEL", "DEBUG"),

		MQTTAddr:   config.String(pref+"MQTT_ADDR", "tcp://127.0.0.1:1883"),
		MQTTUser:   config.String(pref+"MQTT_USER", "mqtt-ingestor"),
		MQTTPass:   config.String(pref+"MQTT_PASS", "notasecurepassword"),
		MQTTShared: config.Bool(pref+"MQTT_SHARED", true),

		NSQPubAddr:  config.String(pref+"NSQ_PUB_ADDR", "127.0.0.1:4150"),
		NSQPubTopic: config.String(pref+"NSQ_PUB_TOPIC", "ValidatorIn"),
		Concurrency: config.Int(pref+"CONCURRENCY", 5),
	}
}
