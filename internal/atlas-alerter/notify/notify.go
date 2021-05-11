package notify

import "github.com/thingspect/atlas/pkg/cache"

// notify contains methods to send notifications and implements the Notifier
// interface.
type notify struct {
	cache cache.Cacher

	appAPIKey   string
	smsID       string
	smsToken    string
	smsPhone    string
	emailAPIKey string
}

// Verify notify implements Notifier.
var _ Notifier = &notify{}

// New builds a new Notifier and returns it.
func New(cache cache.Cacher, appAPIKey, smsID, smsToken, smsPhone,
	emailAPIKey string) Notifier {
	return &notify{
		cache: cache,

		appAPIKey:   appAPIKey,
		smsID:       smsID,
		smsToken:    smsToken,
		smsPhone:    smsPhone,
		emailAPIKey: emailAPIKey,
	}
}
