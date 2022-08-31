package user

import "fmt"

type NotFoundError struct{}

func (uae *NotFoundError) Error() string {
	return fmt.Sprintf("invalid login or password.")
}
