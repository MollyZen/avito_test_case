package dto

type User struct {
	UserID int64 `json:"userID,string" validate:"gte=0,required"`
}
