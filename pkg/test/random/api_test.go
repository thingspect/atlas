// +build !integration

package random

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOrg(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)

			o1 := Org(prefix)
			o2 := Org(prefix)
			t.Logf("o1, o2: %+v, %+v", o1, o2)

			require.True(t, strings.HasPrefix(o1.Name, prefix))
			require.True(t, strings.HasPrefix(o2.Name, prefix))
			require.NotEqual(t, o1, o2)
		})
	}
}

func TestUser(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			u1 := User(prefix, orgID)
			u2 := User(prefix, orgID)
			t.Logf("u1, u2: %+v, %+v", u1, u2)

			require.True(t, strings.HasPrefix(u1.Email, prefix))
			require.True(t, strings.HasPrefix(u2.Email, prefix))
			require.NotEqual(t, u1, u2)
		})
	}
}

func TestDevice(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			d1 := Device(prefix, orgID)
			d2 := Device(prefix, orgID)
			t.Logf("d1, d2: %+v, %+v", d1, d2)

			require.True(t, strings.HasPrefix(d1.UniqId, prefix))
			require.True(t, strings.HasPrefix(d2.UniqId, prefix))
			require.NotEqual(t, d1, d2)
		})
	}
}
