//go:build !integration

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewMemory(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	t.Logf("mem: %+v", mem)
	require.NotNil(t, mem)
}

func TestMemorySetGet(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetGet-" + random.String(10)
	val := random.String(10)

	require.NoError(t, mem.Set(t.Context(), key, val))

	ok, res, err := mem.Get(t.Context(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = mem.Get(t.Context(), "testMemorySetGet-"+
		random.String(10))
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	require.NoError(t, mem.Set(t.Context(), key, random.Bytes(10)))

	ok, res, err = mem.Get(t.Context(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.Equal(t, ErrWrongType, err)
}

func TestMemorySetTTLGetB(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetTTLGetB-" + random.String(10)
	val := random.Bytes(10)

	require.NoError(t, mem.SetTTL(t.Context(), key, val,
		2*time.Second))

	ok, res, err := mem.GetB(t.Context(), key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = mem.GetB(t.Context(), "testMemorySetTTLGetB-"+
		random.String(10))
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	require.NoError(t, mem.SetTTL(t.Context(), key, random.String(10),
		2*time.Second))

	ok, res, err = mem.GetB(t.Context(), key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.Equal(t, ErrWrongType, err)
}

func TestMemorySetTTLGetBShort(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetTTLGetBShort-" + random.String(10)
	val := random.Bytes(10)

	require.NoError(t, mem.SetTTL(t.Context(), key, val,
		time.Millisecond))

	time.Sleep(10 * time.Millisecond)
	ok, res, err := mem.GetB(t.Context(), key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)
}

func TestMemorySetGetI(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetGetI-" + random.String(10)
	val := int64(random.Intn(999))

	require.NoError(t, mem.Set(t.Context(), key, val))

	ok, res, err := mem.GetI(t.Context(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = mem.GetI(t.Context(), "testMemorySetGetI-"+
		random.String(10))
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	require.NoError(t, mem.Set(t.Context(), key, random.String(10)))

	ok, res, err = mem.GetI(t.Context(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.Equal(t, ErrWrongType, err)
}

func TestMemorySetIfNotExist(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetIfNotExist-" + random.String(10)

	ok, err := mem.SetIfNotExist(t.Context(), key, random.Bytes(10))
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = mem.SetIfNotExist(t.Context(), key, random.Bytes(10))
	t.Logf("ok, err: %v, %v", ok, err)
	require.False(t, ok)
	require.NoError(t, err)
}

func TestMemorySetIfNotExistTTL(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetIfNotExistTTL-" + random.String(10)

	ok, err := mem.SetIfNotExistTTL(t.Context(), key, random.Bytes(10),
		2*time.Second)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = mem.SetIfNotExistTTL(t.Context(), key, random.Bytes(10),
		2*time.Second)
	t.Logf("ok, err: %v, %v", ok, err)
	require.False(t, ok)
	require.NoError(t, err)
}

func TestMemorySetIfNotExistTTLShort(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemorySetIfNotExistTTLShort-" + random.String(10)

	ok, err := mem.SetIfNotExistTTL(t.Context(), key, random.Bytes(10),
		time.Millisecond)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	ok, err = mem.SetIfNotExistTTL(t.Context(), key, random.Bytes(10),
		2*time.Second)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)
}

func TestMemoryDel(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	key := "testMemoryDel-" + random.String(10)
	val := random.String(10)

	require.NoError(t, mem.Set(t.Context(), key, val))

	ok, res, err := mem.Get(t.Context(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	err = mem.Del(t.Context(), key)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	ok, res, err = mem.Get(t.Context(), key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)
}

func TestMemoryClose(t *testing.T) {
	t.Parallel()

	mem := NewMemory()
	require.NoError(t, mem.Close())
}
