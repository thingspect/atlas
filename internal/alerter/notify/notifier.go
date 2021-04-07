// Package notify provides functions to send notifications.
package notify

//go:generate mockgen -source notifier.go -destination mock_notifier.go -package notify

import "context"

// Notifier defines the methods provided by a Notify.
type Notifier interface {
	// App sends a push notification to a mobile application. This operation can
	// block based on rate limiting.
	App(ctx context.Context, userKey, subject, body string) error
	// SMS sends an SMS notification. Subjects and bodies will be concatenated
	// with ' - '. This operation can block based on rate limiting.
	SMS(ctx context.Context, phone, subject, body string) error
	// Email sends an email notification. The provider domain used for sending
	// is derived from the organization's email address: "mg." followed by the
	// domain name that follows '@' in the address. This operation can block
	// based on rate limiting.
	Email(ctx context.Context, orgDisplayName, orgEmail, userEmail, subject,
		body string) error
}
