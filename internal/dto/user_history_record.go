package dto

import "time"

type UserHistoryRecord struct {
	UserID    int64     `json:"userID"`
	Segment   string    `json:"segment"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}
