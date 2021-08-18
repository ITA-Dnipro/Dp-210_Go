package timeutils

import "time"

type TimeInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func (i *TimeInterval) Overlaps(other *TimeInterval) bool {
	return max(i.Start, other.Start).Before(min(i.End, other.End))
}
