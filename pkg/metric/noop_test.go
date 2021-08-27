//go:build !integration

package metric

import (
	"fmt"
	"testing"
	"time"

	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNoOp(t *testing.T) {
	t.Parallel()

	metricer := &noOpMetric{}
	t.Logf("metricer: %#v", metricer)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can send %v", lTest), func(t *testing.T) {
			t.Parallel()

			metricer.Incr(random.String(10), nil)
			metricer.Count(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			metricer.Set(random.String(10), random.Intn(99),
				map[string]string{random.String(10): random.String(10)})
			metricer.Timing(random.String(10),
				time.Duration(random.Intn(99))*time.Millisecond, nil)
		})
	}
}
