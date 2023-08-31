package dto

type UserSegmentChange struct {
	UserID        int64          `json:"userID,string" validate:"required"`
	SegmentAdd    []SegmentToAdd `json:"segmentAdd"`
	SegmentRemove []string       `json:"segmentRemove"`
}
