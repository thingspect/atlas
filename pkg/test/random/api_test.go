//go:build !integration

package random

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOrg(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)

			o1 := Org(prefix)
			o2 := Org(prefix)
			t.Logf("o1, o2: %+v, %+v", o1, o2)

			require.NotEqual(t, o1, o2)
			require.True(t, strings.HasPrefix(o1.GetName(), prefix))
			require.True(t, strings.HasPrefix(o2.GetName(), prefix))
			require.True(t, strings.HasPrefix(o1.GetDisplayName(), prefix))
			require.True(t, strings.HasPrefix(o2.GetDisplayName(), prefix))
			require.True(t, strings.HasPrefix(o1.GetEmail(), prefix))
			require.True(t, strings.HasPrefix(o2.GetEmail(), prefix))
		})
	}
}

func TestDevice(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			d1 := Device(prefix, orgID)
			d2 := Device(prefix, orgID)
			t.Logf("d1, d2: %+v, %+v", d1, d2)

			require.NotEqual(t, d1, d2)
			require.True(t, strings.HasPrefix(d1.GetUniqId(), prefix))
			require.True(t, strings.HasPrefix(d2.GetUniqId(), prefix))
			require.True(t, strings.HasPrefix(d1.GetName(), prefix))
			require.True(t, strings.HasPrefix(d2.GetName(), prefix))
			require.NotEmpty(t, d1.GetTags())
			require.NotEmpty(t, d2.GetTags())
		})
	}
}

func TestRule(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			r1 := Rule(prefix, orgID)
			r2 := Rule(prefix, orgID)
			t.Logf("r1, r2: %+v, %+v", r1, r2)

			require.NotEqual(t, r1, r2)
			require.True(t, strings.HasPrefix(r1.GetName(), prefix))
			require.True(t, strings.HasPrefix(r2.GetName(), prefix))
			require.True(t, strings.HasPrefix(r1.GetDeviceTag(), prefix))
			require.True(t, strings.HasPrefix(r2.GetDeviceTag(), prefix))
		})
	}
}

func TestEvent(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			e1 := Event(prefix, orgID)
			e2 := Event(prefix, orgID)
			t.Logf("e1, e2: %+v, %+v", e1, e2)

			require.NotEqual(t, e1, e2)
			require.True(t, strings.HasPrefix(e1.GetUniqId(), prefix))
			require.True(t, strings.HasPrefix(e2.GetUniqId(), prefix))
			require.WithinDuration(t, time.Now(), e1.GetCreatedAt().AsTime(),
				2*time.Second)
			require.WithinDuration(t, time.Now(), e2.GetCreatedAt().AsTime(),
				2*time.Second)
		})
	}
}

func TestAlarm(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()
			ruleID := uuid.NewString()

			a1 := Alarm(prefix, orgID, ruleID)
			a2 := Alarm(prefix, orgID, ruleID)
			t.Logf("a1, a2: %+v, %+v", a1, a2)

			require.NotEqual(t, a1, a2)
			require.True(t, strings.HasPrefix(a1.GetName(), prefix))
			require.True(t, strings.HasPrefix(a2.GetName(), prefix))
			require.NotEmpty(t, a1.GetUserTags())
			require.NotEmpty(t, a2.GetUserTags())
			require.GreaterOrEqual(t, a1.GetRepeatInterval(), int32(1))
			require.GreaterOrEqual(t, a2.GetRepeatInterval(), int32(1))
		})
	}
}

func TestAlert(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			a1 := Alert(prefix, orgID)
			a2 := Alert(prefix, orgID)
			t.Logf("a1, a2: %+v, %+v", a1, a2)

			require.NotEqual(t, a1, a2)
			require.True(t, strings.HasPrefix(a1.GetUniqId(), prefix))
			require.True(t, strings.HasPrefix(a2.GetUniqId(), prefix))
		})
	}
}

func TestUser(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			u1 := User(prefix, orgID)
			u2 := User(prefix, orgID)
			t.Logf("u1, u2: %+v, %+v", u1, u2)

			require.NotEqual(t, u1, u2)
			require.True(t, strings.HasPrefix(u1.GetName(), prefix))
			require.True(t, strings.HasPrefix(u2.GetName(), prefix))
			require.True(t, strings.HasPrefix(u1.GetEmail(), prefix))
			require.True(t, strings.HasPrefix(u2.GetEmail(), prefix))
			require.NotEmpty(t, u1.GetTags())
			require.NotEmpty(t, u2.GetTags())
			require.Empty(t, u1.GetPhone())
			require.Empty(t, u2.GetPhone())
			require.Empty(t, u1.GetAppKey())
			require.Empty(t, u2.GetAppKey())
		})
	}
}

func TestSMSUser(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			u1 := SMSUser(prefix, orgID)
			u2 := SMSUser(prefix, orgID)
			t.Logf("u1, u2: %+v, %+v", u1, u2)

			require.NotEqual(t, u1, u2)
			require.True(t, strings.HasPrefix(u1.GetName(), prefix))
			require.True(t, strings.HasPrefix(u2.GetName(), prefix))
			require.True(t, strings.HasPrefix(u1.GetEmail(), prefix))
			require.True(t, strings.HasPrefix(u2.GetEmail(), prefix))
			require.NotEmpty(t, u1.GetTags())
			require.NotEmpty(t, u2.GetTags())
			require.NotEmpty(t, u1.GetPhone())
			require.NotEmpty(t, u2.GetPhone())
			require.Empty(t, u1.GetAppKey())
			require.Empty(t, u2.GetAppKey())
		})
	}
}

func TestAppUser(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			u1 := AppUser(prefix, orgID)
			u2 := AppUser(prefix, orgID)
			t.Logf("u1, u2: %+v, %+v", u1, u2)

			require.NotEqual(t, u1, u2)
			require.True(t, strings.HasPrefix(u1.GetName(), prefix))
			require.True(t, strings.HasPrefix(u2.GetName(), prefix))
			require.True(t, strings.HasPrefix(u1.GetEmail(), prefix))
			require.True(t, strings.HasPrefix(u2.GetEmail(), prefix))
			require.NotEmpty(t, u1.GetTags())
			require.NotEmpty(t, u2.GetTags())
			require.Empty(t, u1.GetPhone())
			require.Empty(t, u2.GetPhone())
			require.NotEmpty(t, u1.GetAppKey())
			require.NotEmpty(t, u2.GetAppKey())
		})
	}
}

func TestKey(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			prefix := String(10)
			orgID := uuid.NewString()

			k1 := Key(prefix, orgID)
			k2 := Key(prefix, orgID)
			t.Logf("k1, k2: %+v, %+v", k1, k2)

			require.NotEqual(t, k1, k2)
			require.True(t, strings.HasPrefix(k1.GetName(), prefix))
			require.True(t, strings.HasPrefix(k2.GetName(), prefix))
		})
	}
}

func TestTags(t *testing.T) {
	t.Parallel()

	for i := 5; i < 15; i++ {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			t1 := Tags(String(10), i)
			t2 := Tags(String(10), i)
			t.Logf("t1, t2: %v, %v", t1, t2)

			require.Len(t, t1, i)
			require.Len(t, t2, i)
			require.NotEqual(t, t1, t2)
		})
	}
}
