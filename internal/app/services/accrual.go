package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/phrasetagg/gofermart/internal/app/models"
	"github.com/phrasetagg/gofermart/internal/app/repositories"
	"io"
	"net/http"
	"sync"
)

type Accrual struct {
	accrualAddr     string
	client          http.Client
	orderRepository *repositories.Order
}

func NewAccrualService(accrualAddr string, orderRepository *repositories.Order) *Accrual {
	return &Accrual{
		accrualAddr:     accrualAddr,
		client:          http.Client{},
		orderRepository: orderRepository,
	}
}

func (a *Accrual) StartOrderStatusesUpdating() {
	var wg sync.WaitGroup

	if len(a.accrualAddr) == 0 {
		fmt.Println(errors.New("undefined accrual system address"))
		return
	}

	for {
		orders, err := a.orderRepository.GetUnprocessedOrders()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if len(orders) == 0 {
			fmt.Println("No unprocessed orders")
			continue
		}

		wg.Add(len(orders)) // инкрементируем счётчик, сколько горутин нужно подождать

		for _, order := range orders {
			go func(order models.Order) {
				orderInfo, err := a.GetOrderInfo(order.Number)
				if err != nil {
					fmt.Println("Failed to request order status from accrual system" + err.Error())
					wg.Done()
					return
				}

				err = a.orderRepository.ProcessOrderAccrual(orderInfo.Order, orderInfo.Status, orderInfo.Accrual)
				if err != nil {
					fmt.Println("Failed to update order status" + err.Error())
					wg.Done()
					return
				}
				wg.Done()
			}(order)
		}

		wg.Wait() // ждём все горутины
	}
}

func (a *Accrual) GetOrderInfo(orderNumber string) (OrderInfo, error) {
	orderInfo := OrderInfo{}

	response, err := a.client.Get(a.accrualAddr + "/api/orders/" + orderNumber)
	if err != nil {
		return orderInfo, err
	}

	if response.StatusCode != http.StatusOK {
		return orderInfo, errors.New(fmt.Sprintf("%d %s", response.StatusCode, response.Status))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	b, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(b, &orderInfo)

	if err != nil {
		return orderInfo, err
	}

	return orderInfo, nil
}

type OrderInfo struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
