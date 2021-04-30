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

func TestMemorySetGet(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetGet-" + random.String(10)
	val := random.String(10)

	require.NoError(t, mem.Set(context.Background(), key, val))

	ok, res, err := mem.Get(context.Background(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = mem.Get(context.Background(), "testMemorySetGet-"+
		random.String(10))
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	// Test next method.
	ok, resB, err := mem.GetB(context.Background(), key)
	t.Logf("ok, resB, err: %v, %x, %v", ok, resB, err)
	require.False(t, ok)
	require.Empty(t, resB)
	require.Equal(t, errWrongType, err)
}

func TestMemorySetTTLGetB(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetTTLGetB-" + random.String(10)
	val := random.Bytes(10)

	require.NoError(t, mem.SetTTL(context.Background(), key, val,
		2*time.Second))

	ok, res, err := mem.GetB(context.Background(), key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = mem.GetB(context.Background(), "testMemorySetTTLGetB-"+
		random.String(10))
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	// Test next method.
	ok, resI, err := mem.GetI(context.Background(), key)
	t.Logf("ok, resI, err: %v, %x, %v", ok, resI, err)
	require.False(t, ok)
	require.Empty(t, resI)
	require.Equal(t, errWrongType, err)
}

func TestMemorySetTTLGetBShort(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetTTLGetBShort-" + random.String(10)
	val := random.Bytes(10)

	require.NoError(t, mem.SetTTL(context.Background(), key, val,
		time.Millisecond))

	time.Sleep(100 * time.Millisecond)
	ok, res, err := mem.GetB(context.Background(), key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)
}

func TestMemorySetGetI(t *testing.T) {
	t.Parallel()

	mem, err := NewMemory()
	t.Logf("mem, err: %+v, %v", mem, err)
	require.NoError(t, err)

	key := "testMemorySetGetI-" + random.String(10)
	val := int64(random.Intn(999))

	require.NoError(t, mem.SetTTL(context.Background(), key, val,
		2*time.Second))

	ok, res, err := mem.GetI(context.Background(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = mem.GetI(context.Background(), "testMemorySetGetI-"+
		random.String(10))
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	// Test next method.
	ok, resS, err := mem.Get(context.Background(), key)
	t.Logf("ok, resS, err: %v, %v, %v", ok, resS, err)
	require.False(t, ok)
	require.Empty(t, resS)
	require.Equal(t, errWrongType, err)
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
