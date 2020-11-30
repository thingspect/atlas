// Package org provides functions to query and modify organizations in the
// database.
package org

import "time"

// Org represents an organization as stored in the database.
type Org struct {
	ID        string    // id
	Name      string    // name
	CreatedAt time.Time // created_at
	UpdatedAt time.Time // updated_at
}
