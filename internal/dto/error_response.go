package dto

import "fmt"

type ErrorResponse struct {
	Error string `json:"message"`
}

func NewErrorResponse(message interface{}, args ...interface{}) ErrorResponse {
	var v string
	switch t := message.(type) {
	case error:
		v = t.Error()
	case string:
		v = t
	default:
		v = fmt.Sprintf("Unknown type %v", message)
	}
	return ErrorResponse{
		Error: fmt.Sprintf(v, args...),
	}
}
