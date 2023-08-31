package dto

type UserSegmentGet struct {
	UserID       int64          `json:"userID"`
	SegmentAdded []SegmentToAdd `json:"segmentAdd"`
}
