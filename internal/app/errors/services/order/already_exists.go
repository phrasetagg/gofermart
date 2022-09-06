package order

import "fmt"

type AlreadyExistsError struct {
	OrderNumber string
}

func (uae *AlreadyExistsError) Error() string {
	return fmt.Sprintf("you have already uploaded this order number: %s.", uae.OrderNumber)
}

type AlreadyExistsByAnotherUserError struct {
	OrderNumber string
}

func (uaebau *AlreadyExistsByAnotherUserError) Error() string {
	return fmt.Sprintf("another user has already uploaded this order number: %s.", uaebau.OrderNumber)
}
