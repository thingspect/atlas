package metric

import "time"

// noOpMetric mocks sending metrics and implements the Metricer interface.
type noOpMetric struct{}

// Verify noOpMetric implements Metricer.
var _ Metricer = &noOpMetric{}

// Incr mocks incrementing a statsd count metric by 1.
func (n *noOpMetric) Incr(stat string, tags map[string]string) {}

// Count mocks incrementing a statsd count metric.
func (n *noOpMetric) Count(stat string, value int, tags map[string]string) {}

// Set mocks setting a statsd gauge metric. A gauge maintains its value until it
// is next set.
func (n *noOpMetric) Set(stat string, value int, tags map[string]string) {}

// Timing mocks submitting a statsd timing metric in milliseconds.
func (n *noOpMetric) Timing(stat string, value time.Duration,
	tags map[string]string) {
}
