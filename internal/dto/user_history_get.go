package dto

type UserHistoryGet struct {
	User
	Year  int32 `json:"year,string" validate:"required" validation:"gte=0,required"`
	Month int32 `json:"month,string" validate:"required" minimum:"1" maximum:"12" validation:"gte=1,lte=12,required"`
}
