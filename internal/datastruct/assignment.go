package datastruct

import "time"

type Assignment struct {
	ID        int64
	UserID    int64
	SegmentID int64
	UntilDate time.Time
}
