package services

import (
	"errors"
	orderErrors "github.com/phrasetagg/gofermart/internal/app/errors/services/order"
	"github.com/phrasetagg/gofermart/internal/app/models"
	"github.com/phrasetagg/gofermart/internal/app/repositories"
	"regexp"
	"sort"
	"strconv"
)

type Order struct {
	orderRepository *repositories.Order
}

func NewOrderService(orderRepository *repositories.Order) *Order {
	return &Order{orderRepository: orderRepository}
}

type byUploadedAt []models.Order

func (s byUploadedAt) Len() int           { return len(s) }
func (s byUploadedAt) Less(i, j int) bool { return s[j].UploadedAt.After(s[i].UploadedAt) }
func (s byUploadedAt) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (o *Order) GetUserOrders(userID int64) ([]models.Order, error) {
	orders, err := o.orderRepository.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	sort.Sort(byUploadedAt(orders))

	return orders, err
}

func (o *Order) Upload(userID int64, orderNumber string) error {
	order, err := o.orderRepository.GetOrderByNumber(orderNumber)

	// Заказ уже загружен другим пользователем.
	if order.UserID != 0 && order.UserID != userID {
		return &orderErrors.AlreadyExistsByAnotherUserError{OrderNumber: orderNumber}
	}

	// Заказ уже загружен текущим пользователем.
	if order.UserID != 0 && order.UserID == userID {
		return &orderErrors.AlreadyExistsError{OrderNumber: orderNumber}
	}

	// Если возникла ошибка, отличная от отсутствия строк в БД, то возвращаем ее.
	var onfe *orderErrors.NotFoundError
	if err != nil && !errors.As(err, &onfe) {
		return err
	}

	err = o.orderRepository.Create(userID, orderNumber)

	return err
}

// IsValidOrderNumber возвращает true, если номер заказ корректный, иначе false.
func IsValidOrderNumber(orderNumber string) bool {
	re := regexp.MustCompile("^([0-9])+$")

	intNumber, err := strconv.Atoi(orderNumber)
	if err != nil {
		return false
	}

	return re.MatchString(orderNumber) && ValidateLunaAlgorithm(intNumber)
}

// IsNotValidOrderNumber инвертированный метод IsValidOrderNumber.
func IsNotValidOrderNumber(orderNumber string) bool {
	return !IsValidOrderNumber(orderNumber)
}

func ValidateLunaAlgorithm(number int) bool {
	return (number%10+lunaChecksum(number/10))%10 == 0
}

func lunaChecksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
