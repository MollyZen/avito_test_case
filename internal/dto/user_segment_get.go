package dto

type SegmentGet struct {
	UserID       int64    `json:"user-id"`
	SegmentAdded []string `json:"segments"`
}
