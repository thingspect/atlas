// +build !integration

package notify

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewFake(t *testing.T) {
	t.Parallel()

	notifier := NewFake()
	t.Logf("notifier: %#v", notifier)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can notify %v", lTest), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			require.NoError(t, notifier.App(ctx, random.String(10),
				random.String(10), random.String(10)))
			require.NoError(t, notifier.App(ctx, random.String(10),
				random.String(10), random.String(10)))
		})
	}
}
