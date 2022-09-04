package user

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

func GetBalance(userService *services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		user := helpers.GetUserFromCtx(r.Context())
		balance, err := userService.GetBalance(user.ID)

		// 500
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			return
		}

		// 200
		response, err := json.Marshal(balance)
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(response)
		if err != nil {
			return
		}
	}
}

func GetWithdrawals(userService *services.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		user := helpers.GetUserFromCtx(r.Context())
		withdrawals, err := userService.GetWithdrawals(user.ID)

		// 500
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			return
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 200
		response, err := json.Marshal(withdrawals)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			return
		}
	}
}

func RegisterWithDraw(userService *services.User, orderService *services.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		w.Header().Set("content-type", "application/json")

		var request struct {
			OrderNumber string  `json:"order"`
			Sum         float64 `json:"sum"`
		}

		b, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(b, &request)

		// 400
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(`{"error":"invalid request body."}`))
			return
		}

		// 422
		if services.IsNotValidOrderNumber(request.OrderNumber) {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(`{"error":"invalid order number."}`))
			return
		}

		user := helpers.GetUserFromCtx(r.Context())
		balance, err := userService.GetBalance(user.ID)

		// 500
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			return
		}

		// 402
		if balance.Current < request.Sum {
			w.WriteHeader(http.StatusPaymentRequired)
			_, err = w.Write([]byte(`{"error":"not enough funds on the balance."}`))
			return
		}

		err = orderService.Upload(user.ID, request.OrderNumber)

		// 409
		var oae *orderErrors.AlreadyExistsError
		var oaebau *orderErrors.AlreadyExistsByAnotherUserError
		if errors.As(err, &oae) || errors.As(err, &oaebau) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				return
			}
			return
		}

		err = userService.RegisterWithdraw(user.ID, request.OrderNumber, request.Sum)

		// 500
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong."}`))
			return
		}

		// 200
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte(`{"message":"you have successfully registered withdraw."}`))
		if err != nil {
			return
		}
	}
}
