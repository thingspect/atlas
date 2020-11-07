// +build !integration

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestString(t *testing.T) {
	t.Parallel()

	envKey := random.String(10)
	envVal := random.String(10)
	require.NoError(t, os.Setenv(envKey, envVal), "Should set env correctly")
	// Do not unset due to the mechanics of t.Parallel().

	tests := []struct {
		inpKey string
		inpDef string
		res    string
	}{
		{envKey, "", envVal},
		{random.String(10), "default", "default"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := String(lTest.inpKey, lTest.inpDef)
			t.Logf("res: %v", res)
			require.Equal(t, lTest.res, res, "Should match result")
		})
	}
}

func TestStringSlice(t *testing.T) {
	t.Parallel()

	envKey := random.String(10)
	envVal := fmt.Sprintf("%s,%s,%s", random.String(10), random.String(10),
		random.String(10))
	require.NoError(t, os.Setenv(envKey, envVal), "Should set env correctly")

	envKeyNoDelim := random.String(10)
	envValNoDelim := random.String(10)
	require.NoError(t, os.Setenv(envKeyNoDelim, envValNoDelim),
		"Should set env correctly")
	// Do not unset due to the mechanics of t.Parallel().

	tests := []struct {
		inpKey string
		inpDef []string
		res    []string
	}{
		{envKey, []string{}, strings.Split(envVal, ",")},
		{envKeyNoDelim, []string{}, []string{envValNoDelim}},
		{random.String(10), []string{"default"}, []string{"default"}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := StringSlice(lTest.inpKey, lTest.inpDef)
			t.Logf("res: %#v", res)
			require.Equal(t, lTest.res, res, "Should match result")
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	envKey := random.String(10)
	envVal := random.Intn(999)
	require.NoError(t, os.Setenv(envKey, strconv.Itoa(envVal)),
		"Should set env correctly")
	// Do not unset due to the mechanics of t.Parallel().

	tests := []struct {
		inpKey string
		inpDef int
		res    int
	}{
		{envKey, 0, envVal},
		{random.String(10), 99, 99},
		// Do not test conversion failure from env due to use of log.Fatalf().
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := Int(lTest.inpKey, lTest.inpDef)
			t.Logf("res: %v", res)
			require.Equal(t, lTest.res, res, "Should match result")
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	envKey := random.String(10)
	require.NoError(t, os.Setenv(envKey, "true"), "Should set env correctly")
	// Do not unset due to the mechanics of t.Parallel().

	tests := []struct {
		inpKey string
		inpDef bool
		res    bool
	}{
		{envKey, false, true},
		{random.String(10), true, true},
		// Do not test conversion failure from env due to use of log.Fatalf().
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := Bool(lTest.inpKey, lTest.inpDef)
			t.Logf("res: %v", res)
			require.Equal(t, lTest.res, res, "Should match result")
		})
	}
}

func TestDuration(t *testing.T) {
	t.Parallel()

	envKey := random.String(10)
	envVal := random.Intn(999)
	require.NoError(t, os.Setenv(envKey, strconv.Itoa(envVal)+"s"),
		"Should set env correctly")
	// Do not unset due to the mechanics of t.Parallel().

	tests := []struct {
		inpKey string
		inpDef time.Duration
		res    time.Duration
	}{
		{envKey, 0, time.Duration(envVal) * time.Second},
		{random.String(10), time.Minute, time.Minute},
		// Do not test conversion failure from env due to use of log.Fatalf().
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := Duration(lTest.inpKey, lTest.inpDef)
			t.Logf("res: %#v", res)
			require.Equal(t, lTest.res, res, "Should match result")
		})
	}
}
