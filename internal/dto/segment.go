package dto

type Segment struct {
	Slug      string  `json:"slug" minLength:"4" validate:"min=4,required"`
	Percent   float64 `json:"percent,string" validate:"gte=0,lte=100"`
	UntilDate string  `json:"untilDate"`
}
