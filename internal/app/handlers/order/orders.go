package order

import (
	"encoding/json"
	"errors"
	"fmt"
	orderErrors "github.com/phrasetagg/gofermart/internal/app/errors/services/order"
	"github.com/phrasetagg/gofermart/internal/app/helpers"
	"github.com/phrasetagg/gofermart/internal/app/services"
	"io"
	"net/http"
)

func Get(orderService *services.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		w.Header().Set("content-type", "application/json")

		user := helpers.GetUserFromCtx(r.Context())
		orders, err := orderService.GetUserOrders(user.ID)

		// 500
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			if err != nil {
				return
			}
			return
		}

		// 204
		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 200
		response, err := json.Marshal(orders)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			return
		}
	}
}

func Upload(orderService *services.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		w.Header().Set("content-type", "application/json")

		b, err := io.ReadAll(r.Body)

		// 500
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			if err != nil {
				return
			}
			return
		}

		orderNumber := string(b)

		fmt.Println("ORDER NUMBER: " + orderNumber)

		// 400
		if orderNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(`{"error":"invalid request body."}`))
			if err != nil {
				return
			}
			return
		}

		// 422
		if services.IsNotValidOrderNumber(orderNumber) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err = w.Write([]byte(`{"error":"invalid order number."}`))
			if err != nil {
				return
			}
			return
		}

		user := helpers.GetUserFromCtx(r.Context())
		err = orderService.Upload(user.ID, orderNumber)

		// 200
		var oae *orderErrors.AlreadyExistsError
		if errors.As(err, &oae) {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
			if err != nil {
				return
			}
			return
		}

		// 409
		var oaebau *orderErrors.AlreadyExistsByAnotherUserError
		if errors.As(err, &oaebau) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				return
			}
			return
		}

		// 500
		var onfe *orderErrors.NotFoundError
		if err != nil && !errors.As(err, &onfe) {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			if err != nil {
				return
			}
			return
		}

		// 202
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte(`{"message":"new order number accepted for processing."}`))
		if err != nil {
			return
		}
	}
}
