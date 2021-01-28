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

// global returns the global metricer, which is thread-safe.
func global() Metricer {
	metricerMu.Lock()
	defer metricerMu.Unlock()

	return metricer
}

// setGlobal sets a new global metricer.
func setGlobal(m Metricer) {
	metricerMu.Lock()
	defer metricerMu.Unlock()

	metricer = m
}

// Incr increments a statsd count metric by 1.
func Incr(stat string, tags map[string]string) {
	global().Incr(stat, tags)
}

// Count increments a statsd count metric.
func Count(stat string, value int, tags map[string]string) {
	global().Count(stat, value, tags)
}

// Timing submits a statsd timing metric in milliseconds.
func Timing(stat string, value time.Duration, tags map[string]string) {
	global().Timing(stat, value, tags)
}
