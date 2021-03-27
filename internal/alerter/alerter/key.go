package alerter

import "fmt"

func Key(orgID, devID, alarmID, userID string) string {
	return fmt.Sprintf("alerter:org:%s:dev:%s:alarm:%s:user:%s", orgID, devID,
		alarmID, userID)
}
