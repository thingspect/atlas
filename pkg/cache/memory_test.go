//go:build !integration

package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewMemory(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[string]()
	t.Logf("mem: %+v", mtx)
	require.NotNil(t, mtx)

	heap := NewHeap[string]()
	t.Logf("heap: %+v", heap)
	require.NotNil(t, heap)
}

func TestMemorySetGetString(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[string]()
	heap := NewHeap[string]()

	key := "testMemorySetGetString-" + random.String(10)
	val := random.String(10)

	for _, mem := range []Cacher[string]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set get %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.Set(t.Context(), key, val))

			res, err := mem.Get(t.Context(), key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			res, err = mem.Get(t.Context(), "testMemorySetGetString-"+
				random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)
		})
	}
}

func TestMemorySetTTLGetBytes(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[[]byte]()
	heap := NewHeap[[]byte]()

	key := "testMemorySetTTLGetBytes-" + random.String(10)
	val := random.Bytes(10)

	for _, mem := range []Cacher[[]byte]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set TTL get %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.SetTTL(t.Context(), key, val,
				2*time.Second))

			res, err := mem.Get(t.Context(), key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			res, err = mem.Get(t.Context(), "testMemorySetTTLGetBytes-"+
				random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)
		})
	}
}

func TestMemorySetTTLGetBytesShort(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[[]byte]()
	heap := NewHeap[[]byte]()

	key := "testMemorySetTTLGetBytesShort-" + random.String(10)
	val := random.Bytes(10)

	for _, mem := range []Cacher[[]byte]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set TTL get %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.SetTTL(t.Context(), key, val,
				time.Millisecond))

			time.Sleep(10 * time.Millisecond)
			res, err := mem.Get(t.Context(), key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)
		})
	}
}

func TestMemorySetGetInt64(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[int64]()
	heap := NewHeap[int64]()

	key := "testMemorySetGetInt64-" + random.String(10)
	val := int64(random.Intn(999))

	for _, mem := range []Cacher[int64]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set get %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.Set(t.Context(), key, val))

			res, err := mem.Get(t.Context(), key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			res, err = mem.Get(t.Context(), "testMemorySetGetInt64-"+
				random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)
		})
	}
}

func TestMemorySetIfNotExistBytes(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[[]byte]()
	heap := NewHeap[[]byte]()

	key := "testMemorySetIfNotExistBytes-" + random.String(10)

	for _, mem := range []Cacher[[]byte]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set if not exist %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.SetIfNotExist(t.Context(), key,
				random.Bytes(10)))

			require.Equal(t, ErrAlreadyExists, mem.SetIfNotExist(t.Context(),
				key, random.Bytes(10)))
		})
	}
}

func TestMemorySetIfNotExistTTLBytes(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[[]byte]()
	heap := NewHeap[[]byte]()

	key := "testMemorySetIfNotExistTTLBytes-" + random.String(10)

	for _, mem := range []Cacher[[]byte]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set if not exist %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.SetIfNotExistTTL(t.Context(), key,
				random.Bytes(10), 2*time.Second))

			require.Equal(t, ErrAlreadyExists, mem.SetIfNotExistTTL(t.Context(),
				key, random.Bytes(10), 2*time.Second))
		})
	}
}

func TestMemorySetIfNotExistTTLBytesShort(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[[]byte]()
	heap := NewHeap[[]byte]()

	key := "testMemorySetIfNotExistTTLBytesShort-" + random.String(10)

	for _, mem := range []Cacher[[]byte]{mtx, heap} {
		t.Run(fmt.Sprintf("Can set if not exist %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.SetIfNotExistTTL(t.Context(), key,
				random.Bytes(10), time.Millisecond))

			time.Sleep(10 * time.Millisecond)
			require.NoError(t, mem.SetIfNotExistTTL(t.Context(), key,
				random.Bytes(10), 2*time.Second))
		})
	}
}

func TestMemoryIncrInt64(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[int64]()
	heap := NewHeap[int64]()

	key := "testMemoryIncrInt64-" + random.String(10)
	val := int64(random.Intn(999))

	for _, mem := range []Cacher[int64]{mtx, heap} {
		t.Run(fmt.Sprintf("Can not incr %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.Set(t.Context(), key, val))

			require.PanicsWithValue(t, "unimplemented", func() {
				_, _ = mem.Incr(t.Context(), key)
			})
		})
	}
}

func TestMemoryDelString(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[string]()
	heap := NewHeap[string]()

	key := "testMemoryDelString-" + random.String(10)
	val := random.String(10)

	for _, mem := range []Cacher[string]{mtx, heap} {
		t.Run(fmt.Sprintf("Can del %+v", mem), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, mem.Set(t.Context(), key, val))

			res, err := mem.Get(t.Context(), key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			err = mem.Del(t.Context(), key)
			t.Logf("err: %v", err)
			require.NoError(t, err)

			res, err = mem.Get(t.Context(), key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)
		})
	}
}

func TestMemoryClose(t *testing.T) {
	t.Parallel()

	mtx := NewMutex[string]()
	require.NoError(t, mtx.Close())

	heap := NewHeap[string]()
	require.NoError(t, heap.Close())
}
