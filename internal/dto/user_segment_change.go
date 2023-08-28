package dto

type UserSegmentChange struct {
	UserID        int64    `json:"user-id"`
	SegmentAdd    []string `json:"segment-add"`
	SegmentRemove []string `json:"segment-remove"`
}
