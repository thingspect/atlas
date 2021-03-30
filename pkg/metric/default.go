package metric

import (
	"sync"
	"time"
)

// Since metricer is global and may be replaced, locking is required.
var (
	metricer   Metricer
	metricerMu sync.Mutex
)

// getDefault returns the default metricer, which is thread-safe. It is named as
// such because 'default' is a keyword.
func getDefault() Metricer {
	metricerMu.Lock()
	defer metricerMu.Unlock()

	return metricer
}

// setDefault sets a new default metricer.
func setDefault(m Metricer) {
	metricerMu.Lock()
	defer metricerMu.Unlock()

	metricer = m
}

// Incr increments a statsd count metric by 1.
func Incr(stat string, tags map[string]string) {
	getDefault().Incr(stat, tags)
}

// Count increments a statsd count metric.
func Count(stat string, value int, tags map[string]string) {
	getDefault().Count(stat, value, tags)
}

// Set sets a statsd gauge metric. A gauge maintains its value until it is next
// set.
func Set(stat string, value int, tags map[string]string) {
	getDefault().Set(stat, value, tags)
}

// Timing submits a statsd timing metric in milliseconds.
func Timing(stat string, value time.Duration, tags map[string]string) {
	getDefault().Timing(stat, value, tags)
}
