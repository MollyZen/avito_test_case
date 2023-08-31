package dto

type SegmentToDelete struct {
	Slug string `json:"slug" validate:"required"`
}
