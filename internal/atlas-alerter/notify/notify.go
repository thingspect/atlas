package notify

import "github.com/thingspect/atlas/pkg/cache"

// notify contains methods to send notifications and implements the Notifier
// interface.
type notify struct {
	cache cache.Cacher

	appAPIKey   string
	smsSID      string
	smsSecret   string
	smsPhone    string
	emailAPIKey string
}

// Verify notify implements Notifier.
var _ Notifier = &notify{}

// New builds a new Notifier and returns it.
func New(cache cache.Cacher, appAPIKey, smsSID, smsSecret, smsPhone,
	emailAPIKey string) Notifier {
	return &notify{
		cache: cache,

		appAPIKey:   appAPIKey,
		smsSID:      smsSID,
		smsSecret:   smsSecret,
		smsPhone:    smsPhone,
		emailAPIKey: emailAPIKey,
	}
}
