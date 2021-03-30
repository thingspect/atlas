package notify

import "github.com/thingspect/atlas/pkg/cache"

// notify contains methods to send notifications and implements the Notifier
// interface.
type notify struct {
	cache cache.Cacher

	appKey string
}

// Verify notify implements Notifier.
var _ Notifier = &notify{}

// New builds a new Notifier and returns it.
func New(cache cache.Cacher, appKey string) Notifier {
	return &notify{
		cache: cache,

		appKey: appKey,
	}
}
