// +build !integration

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewMemory(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NotNil(t, mem)
	require.NoError(t, err)
}

func TestMemorySetIfNotExist(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetIfNotExist-" + random.String(10)

	ok, err := mem.SetIfNotExist(context.Background(), key, random.Bytes(10))
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = mem.SetIfNotExist(context.Background(), key, random.Bytes(10))
	t.Logf("ok, err: %v, %v", ok, err)
	require.False(t, ok)
	require.NoError(t, err)
}

func TestMemorySetIfNotExistTTL(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetIfNotExistTTL-" + random.String(10)

	ok, err := mem.SetIfNotExistTTL(context.Background(), key, random.Bytes(10),
		2*time.Second)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = mem.SetIfNotExistTTL(context.Background(), key, random.Bytes(10),
		2*time.Second)
	t.Logf("ok, err: %v, %v", ok, err)
	require.False(t, ok)
	require.NoError(t, err)
}

func TestMemorySetIfNotExistTTLShort(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetIfNotExistTTLShort-" + random.String(10)

	ok, err := mem.SetIfNotExistTTL(context.Background(), key, random.Bytes(10),
		time.Millisecond)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	ok, err = mem.SetIfNotExistTTL(context.Background(), key, random.Bytes(10),
		2*time.Second)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)
}

func TestMemoryClose(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	require.NoError(t, mem.Close())
}
