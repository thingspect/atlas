//go:build !integration

package key

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDisabled(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can key %v", i), func(t *testing.T) {
			t.Parallel()

			orgID := uuid.NewString()
			keyID := uuid.NewString()

			key := Disabled(orgID, keyID)
			t.Logf("key: %v", key)

			require.Equal(t, fmt.Sprintf("api:disabled:org:%s:key:%s", orgID,
				keyID), key)
			require.Equal(t, key, Disabled(orgID, keyID))
		})
	}
}
