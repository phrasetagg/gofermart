package order

type NotFoundError struct{}

func (uae *NotFoundError) Error() string {
	return "order not found"
}
