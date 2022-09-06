package user

type NotFoundError struct{}

func (uae *NotFoundError) Error() string {
	return "invalid login or password."
}
