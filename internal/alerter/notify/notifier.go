// Package notify provides functions to send notifications.
package notify

//go:generate mockgen -source notifier.go -destination mock_notifier.go -package notify

import "context"

// Notifier defines the methods provided by a Notify.
type Notifier interface {
	// App sends a push notification to a mobile application. This operation can
	// block based on rate limiting.
	App(ctx context.Context, userKey, subject, body string) error
}
