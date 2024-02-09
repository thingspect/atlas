package device

import "fmt"

// devKey returns a cache key by device ID.
func devKey(orgID, devID string) string {
	return fmt.Sprintf("dao:device:org:%s:dev:%s", orgID, devID)
}

// devKeyByUniqID returns a cache key by device UniqID. This key does not limit
// by org ID and should only be read in the service layer.
func devKeyByUniqID(uniqID string) string {
	return "dao:device:uniqid:" + uniqID
}
