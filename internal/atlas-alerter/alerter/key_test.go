//go:build !integration

package alerter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepeatKey(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can key %v", lTest), func(t *testing.T) {
			t.Parallel()

			orgID := uuid.NewString()
			devID := uuid.NewString()
			alarmID := uuid.NewString()
			userID := uuid.NewString()

			key := repeatKey(orgID, devID, alarmID, userID)
			t.Logf("key: %v", key)

			require.Equal(t, fmt.Sprintf("alerter:repeat:org:%s:dev:%s:alarm:"+
				"%s:user:%s", orgID, devID, alarmID, userID), key)
			require.Equal(t, key, repeatKey(orgID, devID, alarmID, userID))
		})
	}
}
