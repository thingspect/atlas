package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/thingspect/atlas/pkg/consterr"
)

// ErrInvalidSMS is returned when a phone number fails validation.
const ErrInvalidSMS consterr.Error = "unknown or unsupported phone number"

const (
	smsKey       = "notify.sms"
	smsRateDelay = 750 * time.Millisecond
)

// VaildateSMS verifies that a phone number is correct and supported for SMS
// usage.
func (n *notify) VaildateSMS(ctx context.Context, phone string) error {
	lookup, err := n.twilio.lookupCarrier(ctx, phone)
	if err != nil {
		return err
	}

	if lookup.Carrier.Type != "mobile" && lookup.Carrier.Type != "voip" {
		return ErrInvalidSMS
	}

	return nil
}

// SMS sends an SMS notification. Subjects and bodies will be concatenated with
// ' - '. This operation can block based on rate limiting.
func (n *notify) SMS(ctx context.Context, phone, subject, body string) error {
	// Truncate to message limits, supporting up to 2 SMS messages:
	// https://www.twilio.com/docs/glossary/what-sms-character-limit
	msg := subject + " - " + body
	if len(msg) > 300 {
		msg = fmt.Sprintf("%s...", msg[:297])
	}

	// Support modified Twilio rate limit of 1 per second, serially. Twilio will
	// queue up to 4 hours worth of messages (14,400), but at the risk of abuse
	// by fraudulent users:
	// https://support.twilio.com/hc/en-us/articles/115002943027-Understanding-Twilio-Rate-Limits-and-Message-Queues
	ok, err := n.cache.SetIfNotExistTTL(ctx, smsKey, 1, smsRateDelay)
	if err != nil {
		return err
	}
	for !ok {
		time.Sleep(smsRateDelay)

		ok, err = n.cache.SetIfNotExistTTL(ctx, smsKey, 1, smsRateDelay)
		if err != nil {
			return err
		}
	}

	return n.twilio.sendSMS(ctx, phone, msg)
}
