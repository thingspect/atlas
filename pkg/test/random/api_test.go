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

			require.NotEqual(t, o1, o2)
			require.True(t, strings.HasPrefix(o1.Name, prefix))
			require.True(t, strings.HasPrefix(o2.Name, prefix))
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

			require.NotEqual(t, d1, d2)
			require.True(t, strings.HasPrefix(d1.UniqId, prefix))
			require.True(t, strings.HasPrefix(d2.UniqId, prefix))
			require.True(t, strings.HasPrefix(d1.Name, prefix))
			require.True(t, strings.HasPrefix(d2.Name, prefix))
			require.GreaterOrEqual(t, len(d1.Tags), 1)
			require.GreaterOrEqual(t, len(d2.Tags), 1)
		})
	}
}

func TestRule(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			r1 := Rule(prefix, orgID)
			r2 := Rule(prefix, orgID)
			t.Logf("r1, r2: %+v, %+v", r1, r2)

			require.NotEqual(t, r1, r2)
			require.True(t, strings.HasPrefix(r1.Name, prefix))
			require.True(t, strings.HasPrefix(r2.Name, prefix))
			require.True(t, strings.HasPrefix(r1.DeviceTag, prefix))
			require.True(t, strings.HasPrefix(r2.DeviceTag, prefix))
		})
	}
}

func TestAlarm(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()
			ruleID := uuid.NewString()

			a1 := Alarm(prefix, orgID, ruleID)
			a2 := Alarm(prefix, orgID, ruleID)
			t.Logf("a1, a2: %+v, %+v", a1, a2)

			require.NotEqual(t, a1, a2)
			require.True(t, strings.HasPrefix(a1.Name, prefix))
			require.True(t, strings.HasPrefix(a2.Name, prefix))
			require.GreaterOrEqual(t, len(a1.UserTags), 1)
			require.GreaterOrEqual(t, len(a2.UserTags), 1)
			require.GreaterOrEqual(t, a1.RepeatInterval, int32(1))
			require.GreaterOrEqual(t, a2.RepeatInterval, int32(1))
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

			require.NotEqual(t, u1, u2)
			require.True(t, strings.HasPrefix(u1.Email, prefix))
			require.True(t, strings.HasPrefix(u2.Email, prefix))
			require.GreaterOrEqual(t, len(u1.Tags), 1)
			require.GreaterOrEqual(t, len(u2.Tags), 1)
		})
	}
}

func TestTags(t *testing.T) {
	t.Parallel()

	for i := 5; i < 15; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			t1 := Tags(String(10), lTest)
			t2 := Tags(String(10), lTest)
			t.Logf("t1, t2: %v, %v", t1, t2)

			require.Len(t, t1, lTest)
			require.Len(t, t2, lTest)
			require.NotEqual(t, t1, t2)
		})
	}
}
