package notify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thingspect/atlas/pkg/cache"
)

const (
	emailKey       string = "notify.email"
	emailRateDelay        = 333 * time.Millisecond
)

// Email sends an email notification. This operation can block based on rate
// limiting.
func (n *notify) Email(
	ctx context.Context, displayName, from, to, subject, body string,
) error {
	// Mailgun does not employ a rate limit, so default to 3 per second,
	// serially.
	err := n.cache.SetIfNotExistTTL(ctx, emailKey, "", emailRateDelay)
	if err != nil && !errors.Is(err, cache.ErrAlreadyExists) {
		return err
	}
	for errors.Is(err, cache.ErrAlreadyExists) {
		time.Sleep(emailRateDelay)

		err = n.cache.SetIfNotExistTTL(ctx, emailKey, "", emailRateDelay)
		if err != nil && !errors.Is(err, cache.ErrAlreadyExists) {
			return err
		}
	}

	sender := fmt.Sprintf("%s <%s>", displayName, from)

	return n.mailgun.sendEmail(ctx, sender, to, subject, body)
}
