package notify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gregdel/pushover"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/metric"
)

// ErrInvalidPushover is returned when a user key fails validation.
const ErrInvalidPushover consterr.Error = "unknown user key"

const (
	appKey       = "notify.app"
	appRateDelay = 500 * time.Millisecond
)

// ValidateApp verifies that a mobile application user key is valid.
func (n *notify) ValidateApp(userKey string) error {
	po := pushover.New(n.appAPIKey)
	recipient := pushover.NewRecipient(userKey)

	// GetRecipientDetails does not return sentinel errors via the API, return
	// ErrInvalidPushover based on status.
	det, err := po.GetRecipientDetails(recipient)
	if errors.Is(err, pushover.ErrInvalidRecipientToken) ||
		(det != nil && det.Status != 1) {
		return ErrInvalidPushover
	}
	if err != nil {
		return err
	}

	return nil
}

// App sends a push notification to a mobile application. This operation can
// block based on rate limiting.
func (n *notify) App(ctx context.Context, userKey, subject, body string) error {
	po := pushover.New(n.appAPIKey)
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

	resp, err := po.SendMessage(msg, recipient)
	// Set remaining message limit if present, regardless of error.
	if resp != nil && resp.Limit != nil {
		metric.Set(appKey+".remaining", resp.Limit.Remaining, nil)
	}

	return err
}
