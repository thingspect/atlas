package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/valkey-io/valkey-go"
)

// valkeyCache contains methods to create and query data in Valkey and
// implements the Cacher interface.
type valkeyCache[V Cacheable] struct {
	client valkey.Client
}

// Verify valkeyCache implements Cacher.
var _ Cacher[string] = &valkeyCache[string]{}

// NewValkey builds and verifies a new Cacher and returns it and an error value.
func NewValkey[V Cacheable](valkeyAddr string) (Cacher[V], error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyAddr},
	})
	if err != nil {
		return nil, err
	}

	return &valkeyCache[V]{client: client}, nil
}

// Set sets key to value.
func (v *valkeyCache[V]) Set(ctx context.Context, key string, value V) error {
	return v.SetTTL(ctx, key, value, 0)
}

// SetTTL sets key to value with expiration.
func (v *valkeyCache[V]) SetTTL(ctx context.Context, key string, value V,
	exp time.Duration,
) error {
	// Switching on - and converting back to - type parameters is not supported.
	// Use JSON bytes as a workaround. This has the added benefit of supporting
	// primitives as understood by Valkey.
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if exp == 0 {
		return v.client.Do(ctx, v.client.B().Set().Key(key).
			Value(valkey.BinaryString(b)).Build()).Error()
	}

	return v.client.Do(ctx, v.client.B().Set().Key(key).
		Value(valkey.BinaryString(b)).Px(exp).Build()).Error()
}

// Get retrieves a value by key. If the key does not exist, ErrNotFound is
// returned.
func (v *valkeyCache[V]) Get(ctx context.Context, key string) (V, error) {
	b, err := v.client.Do(ctx, v.client.B().Get().Key(key).Build()).AsBytes()
	if valkey.IsValkeyNil(err) {
		return *new(V), ErrNotFound
	}
	if err != nil {
		return *new(V), err
	}

	var item V
	if err = json.Unmarshal(b, &item); err != nil {
		return *new(V), err
	}

	return item, nil
}

// SetIfNotExist sets key to value if the key does not exist. If the key already
// exists, ErrAlreadyExists is returned.
func (v *valkeyCache[V]) SetIfNotExist(ctx context.Context, key string,
	value V,
) error {
	return v.SetIfNotExistTTL(ctx, key, value, 0)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If the key already exists, ErrAlreadyExists is returned.
func (v *valkeyCache[V]) SetIfNotExistTTL(ctx context.Context, key string,
	value V, exp time.Duration,
) error {
	// Switching on - and converting back to - type parameters is not supported.
	// Use JSON bytes as a workaround. This has the added benefit of supporting
	// primitives as understood by Valkey.
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var ok bool
	if exp == 0 {
		ok, err = v.client.Do(ctx, v.client.B().Setnx().Key(key).
			Value(valkey.BinaryString(b)).Build()).AsBool()
	} else {
		ok, err = v.client.Do(ctx, v.client.B().Set().Key(key).
			Value(valkey.BinaryString(b)).Nx().Px(exp).Build()).AsBool()
		if valkey.IsValkeyNil(err) {
			// ok is already set to true for a Valkey nil.
			err = nil
		}
	}
	if err != nil {
		return err
	}
	if !ok {
		return ErrAlreadyExists
	}

	return nil
}

// Incr increments an int64 value at key by one. If the key does not exist, the
// value is set to 1. The incremented value is returned.
func (v *valkeyCache[V]) Incr(ctx context.Context, key string) (int64, error) {
	i, err := v.client.Do(ctx, v.client.B().Incr().Key(key).Build()).AsInt64()
	if err != nil {
		return 0, err
	}

	return i, nil
}

// Del removes the specified key. A key is ignored if it does not exist.
func (v *valkeyCache[V]) Del(ctx context.Context, key string) error {
	return v.client.Do(ctx, v.client.B().Del().Key(key).Build()).Error()
}

// Close closes the Cacher, releasing any open resources.
func (v *valkeyCache[V]) Close() error {
	v.client.Close()

	return nil
}
