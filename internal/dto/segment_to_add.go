package dto

import "time"

type SegmentToAdd struct {
	Name      string    `json:"name"`
	UntilDate time.Time `json:"untilDate"`
}
