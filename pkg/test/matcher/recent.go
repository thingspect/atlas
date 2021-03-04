package matcher

import (
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
)

// recentMatcher implements the gomock.Matcher interface.
type recentMatcher struct {
	d time.Duration
}

// Verify recentMatcher implements gomock.Matcher.
var _ gomock.Matcher = &recentMatcher{}

// NewRecentMatcher builds a new gomock.Matcher and returns it.
func NewRecentMatcher(d time.Duration) gomock.Matcher {
	return &recentMatcher{
		d: d,
	}
}

// Matches returns whether x is a match.
func (pm *recentMatcher) Matches(x interface{}) bool {
	t, ok := x.(time.Time)
	if !ok {
		return false
	}

	if (t.Before(time.Now()) && time.Since(t) < pm.d) || (t.After(time.Now()) &&
		time.Until(t) < pm.d) {
		return true
	}

	return false
}

// String describes what the matcher matches.
func (pm *recentMatcher) String() string {
	return fmt.Sprintf("is within %v", pm.d)
}
