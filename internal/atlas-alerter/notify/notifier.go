// Package notify provides functions to send notifications.
package notify

//go:generate mockgen -source notifier.go -destination mock_notifier.go -package notify

import "context"

// Notifier defines the methods provided by a Notify.
type Notifier interface {
	// ValidateApp verifies that a mobile application user key is valid.
	ValidateApp(ctx context.Context, userKey string) error
	// App sends a push notification to a mobile application. This operation can
	// block based on rate limiting.
	App(ctx context.Context, userKey, subject, body string) error
	// VaildateSMS verifies that a phone number is correct and supported for SMS
	// usage.
	VaildateSMS(ctx context.Context, phone string) error
	// SMS sends an SMS notification. Subjects and bodies will be concatenated
	// with ' - '. This operation can block based on rate limiting.
	SMS(ctx context.Context, phone, subject, body string) error
	// Email sends an email notification. This operation can block based on rate
	// limiting.
	Email(ctx context.Context, displayName, from, to, subject,
		body string) error
}
