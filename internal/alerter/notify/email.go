package notify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/thingspect/atlas/pkg/consterr"
)

const (
	errInvalidEmail consterr.Error = "invalid email address"
	emailKey        string         = "notify.email"
)

// Email sends an email notification. The provider domain used for sending is
// derived from the organization's email address: "mg." followed by the domain
// name that follows '@' in the address. This operation can block based on rate
// limiting.
func (n *notify) Email(ctx context.Context, orgDisplayName, orgEmail, userEmail,
	subject, body string) error {
	// Build provider domain.
	domParts := strings.SplitN(orgEmail, "@", 2)
	if len(domParts) != 2 {
		return errInvalidEmail
	}

	client := mailgun.NewMailgun("mg."+domParts[1], n.emailAPIKey)

	// Build message.
	sender := fmt.Sprintf("%s <%s>", orgDisplayName, orgEmail)
	msg := client.NewMessage(sender, subject, body, userEmail)

	// Mailgun does not currently employ a rate limit, so default to 3 per
	// second, serially.
	ok, err := n.cache.SetIfNotExistTTL(ctx, emailKey, 1, 333*time.Millisecond)
	if err != nil {
		return err
	}
	for !ok {
		time.Sleep(333 * time.Millisecond)

		ok, err = n.cache.SetIfNotExistTTL(ctx, emailKey, 1,
			333*time.Millisecond)
		if err != nil {
			return err
		}
	}

	_, _, err = client.Send(ctx, msg)

	return err
}
