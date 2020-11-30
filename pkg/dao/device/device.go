// Package device provides functions to query and modify devices in the
// database.
package device

import "time"

// Device represents a device as stored in the database.
type Device struct {
	ID        string    // id
	OrgID     string    // org_id
	UniqID    string    // uniq_id
	Token     string    // token
	CreatedAt time.Time // created_at
	UpdatedAt time.Time // updated_at
}
