package order

import "fmt"

type NotFoundError struct{}

func (uae *NotFoundError) Error() string {
	return fmt.Sprintf("order not found")
}
