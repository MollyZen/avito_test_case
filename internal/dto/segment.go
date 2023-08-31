package dto

type Segment struct {
	Slug    string  `json:"slug" validate:"required"`
	Percent float64 `json:"percent" validate:"gte=0,lte=100"`
}
