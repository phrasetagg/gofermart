package user

import "fmt"

type AlreadyExistsError struct {
	Login string
}

func (uae *AlreadyExistsError) Error() string {
	return fmt.Sprintf("user with login %s already exists.", uae.Login)
}
