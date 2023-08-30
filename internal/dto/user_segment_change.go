package dto

type UserSegmentChange struct {
	UserID        int64          `json:"userID"`
	SegmentAdd    []SegmentToAdd `json:"segmentAdd"`
	SegmentRemove []string       `json:"segmentRemove"`
}
