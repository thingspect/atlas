package notify

import (
	"context"
	"fmt"
	"time"
)

const (
	emailKey       string = "notify.email"
	emailRateDelay        = 333 * time.Millisecond
)

// Email sends an email notification. This operation can block based on rate
// limiting.
func (n *notify) Email(ctx context.Context, displayName, orgEmail, userEmail,
	subject, body string) error {
	// Mailgun does not employ a rate limit, so default to 3 per second,
	// serially.
	ok, err := n.cache.SetIfNotExistTTL(ctx, emailKey, 1, emailRateDelay)
	if err != nil {
		return err
	}
	for !ok {
		time.Sleep(emailRateDelay)

		ok, err = n.cache.SetIfNotExistTTL(ctx, emailKey, 1, emailRateDelay)
		if err != nil {
			return err
		}
	}

	sender := fmt.Sprintf("%s <%s>", displayName, orgEmail)

	return n.mailgun.sendEmail(ctx, sender, userEmail, subject, body)
}
