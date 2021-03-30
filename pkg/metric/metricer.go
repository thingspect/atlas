// Package metric provides functions to write metrics with a friendly API. An
// interface and implementation was chosen over simpler package wrappers to
// support changing implementations.
package metric

import "time"

// Metricer defines the methods provided by a Metric.
type Metricer interface {
	// Incr increments a statsd count metric by 1.
	Incr(stat string, tags map[string]string)
	// Count increments a statsd count metric.
	Count(stat string, value int, tags map[string]string)
	// Set sets a statsd gauge metric. A gauge maintains its value until it is
	// next set.
	Set(stat string, value int, tags map[string]string)
	// Timing submits a statsd timing metric in milliseconds.
	Timing(stat string, value time.Duration, tags map[string]string)
}
