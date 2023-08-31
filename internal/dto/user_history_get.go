package dto

type UserHistoryGet struct {
	User
	Year  int32 `json:"year,string" validation:"gte=0"`
	Month int32 `json:"month,string" validation:"gte=1,lte=12"`
}
