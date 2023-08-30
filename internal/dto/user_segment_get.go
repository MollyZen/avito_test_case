package dto

type SegmentGet struct {
	UserID       int64    `json:"userID"`
	SegmentAdded []string `json:"segments"`
}
