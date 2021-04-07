package notify

import "github.com/thingspect/atlas/pkg/cache"

// notify contains methods to send notifications and implements the Notifier
// interface.
type notify struct {
	cache cache.Cacher

	appAPIKey     string
	smsAccountSID string
	smsAuthToken  string
	smsPhone      string
	emailAPIKey   string
}

// Verify notify implements Notifier.
var _ Notifier = &notify{}

// New builds a new Notifier and returns it.
func New(cache cache.Cacher, appAPIKey, smsAccountSID, smsAuthToken, smsPhone,
	emailAPIKey string) Notifier {
	return &notify{
		cache: cache,

		appAPIKey:     appAPIKey,
		smsAccountSID: smsAccountSID,
		smsAuthToken:  smsAuthToken,
		smsPhone:      smsPhone,
		emailAPIKey:   emailAPIKey,
	}
}
