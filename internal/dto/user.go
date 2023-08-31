package dto

type User struct {
	UserID int64 `json:"userID,string" validate:"required"`
}
