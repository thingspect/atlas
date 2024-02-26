//go:build !integration

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

	for i := range 5 {
		t.Run(fmt.Sprintf("Can notify %v", i), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			require.NoError(t, notifier.ValidateApp(random.String(10)))
			require.NoError(t, notifier.App(ctx, random.String(10),
				random.String(10), random.String(10)))
			require.NoError(t, notifier.ValidateSMS(ctx, random.String(10)))
			require.NoError(t, notifier.SMS(ctx, random.String(10),
				random.String(10), random.String(10)))
			require.NoError(t, notifier.Email(ctx, random.String(10),
				random.String(10), random.String(10), random.String(10),
				random.String(10)))
		})
	}
}
