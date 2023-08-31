package datastruct

import "time"

type History struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"userid"`
	SegmentID   int64     `db:"segmentid"`
	OperationID int64     `db:"operationid"`
	TimeStamp   time.Time `db:"timestamp"`
}
