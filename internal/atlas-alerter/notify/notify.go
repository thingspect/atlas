package notify

import "github.com/thingspect/atlas/pkg/cache"

// notify contains methods to send notifications and implements the Notifier
// interface.
type notify struct {
	cache cache.Cacher

	appAPIKey string
	twilio    *twilio
	mailgun   *mailgun
}

// Verify notify implements Notifier.
var _ Notifier = &notify{}

// New builds a new Notifier and returns it.
func New(cache cache.Cacher, appAPIKey, smsKeyID, smsAccountID, smsKeySecret,
	smsPhone, emailDomain, emailAPIKey string) Notifier {
	return &notify{
		cache: cache,

		appAPIKey: appAPIKey,
		twilio: &twilio{
			keySID:     smsKeyID,
			accountSID: smsAccountID,
			keySecret:  smsKeySecret,
			phone:      smsPhone,
		},
		mailgun: &mailgun{
			domain: emailDomain,
			apiKey: emailAPIKey,
		},
	}
}
