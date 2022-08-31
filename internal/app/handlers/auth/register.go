package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	userErrors "github.com/phrasetagg/gofermart/internal/app/errors/services/user"
	"github.com/phrasetagg/gofermart/internal/app/services"
	"io"
	"net/http"
)

func Register(userService *services.User, authService *services.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		w.Header().Set("content-type", "application/json")

		var request struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		b, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(b, &request)

		// 400
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		var uae *userErrors.AlreadyExistsError
		err = userService.Register(request.Login, request.Password)

		// 409
		if errors.As(err, &uae) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		// 500
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"error":"something went wrong"}`))
			return
		}

		// Аутентифицируем пользваотеля.
		authToken := authService.GenerateAuthToken(request.Login)
		http.SetCookie(
			w,
			&http.Cookie{
				Name:  services.AuthTokenName,
				Value: authToken,
			})

		// 200
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"message":"you have been registered successfully"}`))
		if err != nil {
			return
		}
	}
}
