package dto

type SegmentToAdd struct {
	Slug      string `json:"slug" minLength:"4" validate:"min=4,required"`
	UntilDate string `json:"untilDate"`
}
