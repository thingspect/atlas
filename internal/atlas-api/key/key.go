// Package key provides functions to generate cache keys.
package key

import "fmt"

// Disabled returns a cache key to support disabled API keys.
func Disabled(orgID, keyID string) string {
	return fmt.Sprintf("api:disabled:org:%s:key:%s", orgID, keyID)
}
