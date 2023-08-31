package dto

type SegmentToDelete struct {
	Slug string `json:"slug" minLength:"4" validate:"min=4,required"`
}
