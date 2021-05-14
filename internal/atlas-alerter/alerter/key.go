package alerter

import "fmt"

// repeatKey returns a cache key to support repeat intervals.
func repeatKey(orgID, devID, alarmID, userID string) string {
	return fmt.Sprintf("alerter:repeat:org:%s:dev:%s:alarm:%s:user:%s", orgID,
		devID, alarmID, userID)
}
