package datastruct

import "time"

type Assignment struct {
	UserID    int64
	SegmentID int64
	UntilDate time.Time
}
