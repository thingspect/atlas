package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/kevinburke/twilio-go"
)

const (
	smsKey       = "notify.sms"
	smsRateDelay = 750 * time.Millisecond
)

// SMS sends an SMS notification. Subjects and bodies will be concatenated with
// ' - '. This operation can block based on rate limiting.
func (n *notify) SMS(ctx context.Context, phone, subject, body string) error {
	client := twilio.NewClient(n.smsID, n.smsToken, nil)

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

	_, err = client.Messages.SendMessage(n.smsPhone, phone, msg, nil)

	return err
}
