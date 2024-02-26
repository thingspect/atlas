//go:build !integration

package device

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestDevKey(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can key %v", i), func(t *testing.T) {
			t.Parallel()

			orgID := uuid.NewString()
			devID := uuid.NewString()

			key := devKey(orgID, devID)
			t.Logf("key: %v", key)

			require.Equal(t, fmt.Sprintf("dao:device:org:%s:dev:%s", orgID,
				devID), key)
			require.Equal(t, key, devKey(orgID, devID))
		})
	}
}

func TestDevKeyByUniqID(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can key %v", i), func(t *testing.T) {
			t.Parallel()

			uniqID := random.String(16)

			key := devKeyByUniqID(uniqID)
			t.Logf("key: %v", key)

			require.Equal(t, "dao:device:uniqid:"+uniqID, key)
			require.Equal(t, key, devKeyByUniqID(uniqID))
		})
	}
}
