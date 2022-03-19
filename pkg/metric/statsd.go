package metric

import (
	"strings"
	"time"

	"github.com/smira/go-statsd"
	"github.com/thingspect/atlas/pkg/alog"
)

// statsD contains methods to send metrics to StatsD and implements the
// Metricer interface.
type statsD struct {
	client *statsd.Client
}

// Verify statsd implements Metricer.
var _ Metricer = &statsD{}

// SetStatsD builds a new StatsD Metric and sets it to the default metricer.
func SetStatsD(addr, prefix string) {
	if addr == "" {
		alog.Error("SetStatsD addr not found, continuing to use noOpMetric")

		return
	}

	if !strings.HasSuffix(prefix, ".") {
		prefix += "."
	}

	setDefault(&statsD{
		client: statsd.NewClient(addr,
			statsd.TagStyle(statsd.TagFormatDatadog),
			statsd.MetricPrefix(prefix)),
	})
}

// Incr increments a statsd count metric by 1.
func (s *statsD) Incr(stat string, tags map[string]string) {
	s.Count(stat, 1, tags)
}

// Count increments a statsd count metric.
func (s *statsD) Count(stat string, value int, tags map[string]string) {
	s.client.Incr(stat, int64(value), tagsToSTags(tags)...)
}

// Set sets a statsd gauge metric. A gauge maintains its value until it is next
// set.
func (s *statsD) Set(stat string, value int, tags map[string]string) {
	s.client.Gauge(stat, int64(value), tagsToSTags(tags)...)
}

// Timing submits a statsd timing metric in milliseconds.
func (s *statsD) Timing(
	stat string, value time.Duration, tags map[string]string,
) {
	s.client.Timing(stat, int64(value/time.Millisecond), tagsToSTags(tags)...)
}

// tagsToSTags converts a tags map to statsd.Tag slice.
func tagsToSTags(tags map[string]string) []statsd.Tag {
	var sTags []statsd.Tag
	if len(tags) > 0 {
		for k, v := range tags {
			sTags = append(sTags, statsd.StringTag(k, v))
		}
	}

	return sTags
}
