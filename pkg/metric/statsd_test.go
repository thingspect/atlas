//go:build !integration

package metric

import (
	"fmt"
	"testing"
	"time"

	"github.com/smira/go-statsd"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestStatsD(t *testing.T) {
	t.Parallel()

	metStats := statsD{
		client: statsd.NewClient("127.0.0.1:8125",
			statsd.TagStyle(statsd.TagFormatGraphite),
			statsd.MetricPrefix("teststatsd.")),
	}
	t.Logf("metStats: %#v", metStats)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can send %v", lTest), func(t *testing.T) {
			t.Parallel()

			metStats.Incr(random.String(10), nil)
			metStats.Count(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			metStats.Set(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			metStats.Timing(random.String(10),
				time.Duration(random.Intn(99))*time.Millisecond, nil)
		})
	}
}

func TestSetStatsD(t *testing.T) {
	SetStatsD("127.0.0.1:8125", "testnewstatsd")

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can send %v", lTest), func(t *testing.T) {
			t.Parallel()

			Incr(random.String(10), nil)
			Count(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			Set(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			Timing(random.String(10),
				time.Duration(random.Intn(99))*time.Millisecond, nil)
		})
	}
}

func TestNewStatsDNoAddr(t *testing.T) {
	t.Parallel()

	SetStatsD("", "testnewstatsdnoaddr")

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can send %v", lTest), func(t *testing.T) {
			t.Parallel()

			Incr(random.String(10), nil)
			Count(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			Timing(random.String(10),
				time.Duration(random.Intn(99))*time.Millisecond, nil)
		})
	}
}
