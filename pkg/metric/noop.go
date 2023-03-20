package metric

import "time"

// noOpMetric mocks sending metrics and implements the Metricer interface.
type noOpMetric struct{}

// Verify noOpMetric implements Metricer.
var _ Metricer = &noOpMetric{}

// Incr mocks incrementing a statsd count metric by 1.
func (n *noOpMetric) Incr(_ string, _ map[string]string) {}

// Count mocks incrementing a statsd count metric.
func (n *noOpMetric) Count(_ string, _ int, _ map[string]string) {}

// Set mocks setting a statsd gauge metric. A gauge maintains its value until it
// is next set.
func (n *noOpMetric) Set(_ string, _ int, _ map[string]string) {}

// Timing mocks submitting a statsd timing metric in milliseconds.
func (n *noOpMetric) Timing(_ string, _ time.Duration, _ map[string]string) {}
