package middlewares

import (
	"context"
	"github.com/phrasetagg/gofermart/internal/app/repositories"
	"github.com/phrasetagg/gofermart/internal/app/services"
	"net/http"
)

type Auth struct {
	authService    *services.Auth
	userRepository *repositories.User
}

func NewAuthMiddleware(authService *services.Auth, userRepository *repositories.User) *Auth {
	return &Auth{
		authService:    authService,
		userRepository: userRepository,
	}
}

type CtxPropName string

const UserCtxPropName CtxPropName = "user"

func (a *Auth) CheckAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authToken := ""
			cookies := r.Cookies()

			for _, cookie := range cookies {
				if cookie.Name == services.AuthTokenName {
					authToken = cookie.Value
				}
			}

			if authToken == "" || !a.authService.ValidateAuthToken(authToken) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userLogin := a.authService.GetUserLoginFromAuthToken(authToken)

			user, err := a.userRepository.GetUserByLogin(userLogin)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), UserCtxPropName, user))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
