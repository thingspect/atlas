//go:build !integration

package alerter

import (
	"fmt"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

func TestRepeatKey(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can key %v", i), func(t *testing.T) {
			t.Parallel()

			orgID := uuid.NewV7().String()
			devID := uuid.NewV7().String()
			alarmID := uuid.NewV7().String()
			userID := uuid.NewV7().String()

			key := repeatKey(orgID, devID, alarmID, userID)
			t.Logf("key: %v", key)

			require.Equal(t, fmt.Sprintf("alerter:repeat:org:%s:dev:%s:alarm:"+
				"%s:user:%s", orgID, devID, alarmID, userID), key)
			require.Equal(t, key, repeatKey(orgID, devID, alarmID, userID))
		})
	}
}
