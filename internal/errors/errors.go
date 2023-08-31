package errors

import "fmt"

type NoSuchUserError struct {
	UserID int64
}

func (n *NoSuchUserError) Error() string {
	return fmt.Sprintf("User with ID %d does not exist", n.UserID)
}
