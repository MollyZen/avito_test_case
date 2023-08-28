package datastruct

import "time"

type History struct {
	ID          int64
	UserID      int64
	SegmentID   int64
	OperationID int64
	TimeStamp   time.Time
}
