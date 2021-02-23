// +build !integration

package lora

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

	loraer := NewFake()
	t.Logf("loraer: %#v", loraer)

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can lora %v", lTest), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			require.NoError(t, loraer.CreateGateway(ctx, random.String(16)))
			require.NoError(t, loraer.DeleteGateway(ctx, random.String(16)))
			require.NoError(t, loraer.CreateDevice(ctx, random.String(16),
				random.String(32)))
			require.NoError(t, loraer.DeleteDevice(ctx, random.String(16)))
		})
	}
}
