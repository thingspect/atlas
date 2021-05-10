package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/gregdel/pushover"
	"github.com/thingspect/atlas/pkg/metric"
)

const (
	appKey       = "notify.app"
	appRateDelay = 500 * time.Millisecond
)

// App sends a push notification to a mobile application. This operation can
// block based on rate limiting.
func (n *notify) App(ctx context.Context, userKey, subject, body string) error {
	app := pushover.New(n.appAPIKey)
	recipient := pushover.NewRecipient(userKey)

	// Truncate to subject and body limits: https://pushover.net/api#limits
	if len(subject) > 250 {
		subject = fmt.Sprintf("%s...", subject[:247])
	}
	if len(body) > 1024 {
		body = fmt.Sprintf("%s...", body[:1024])
	}
	msg := pushover.NewMessageWithTitle(body, subject)

	// Support modified Pushover rate limit of 2 per second, serially:
	// https://pushover.net/api#friendly
	ok, err := n.cache.SetIfNotExistTTL(ctx, appKey, 1, appRateDelay)
	if err != nil {
		return err
	}
	for !ok {
		time.Sleep(appRateDelay)

		ok, err = n.cache.SetIfNotExistTTL(ctx, appKey, 1, appRateDelay)
		if err != nil {
			return err
		}
	}

	resp, err := app.SendMessage(msg, recipient)
	// Set remaining message limit if present, regardless of error.
	if resp != nil && resp.Limit != nil {
		metric.Set(appKey+".remaining", resp.Limit.Remaining, nil)
	}

	return err
}
