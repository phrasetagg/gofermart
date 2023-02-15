package helpers

import (
	"context"
	"github.com/phrasetagg/gofermart/internal/app/middlewares"
	userModels "github.com/phrasetagg/gofermart/internal/app/models/user"
)

// GetUserFromCtx возвращает объект пользователя из контекста.
func GetUserFromCtx(ctx context.Context) *userModels.User {
	var user *userModels.User

	rawUser := ctx.Value(middlewares.UserCtxPropName)

	if rawUser == nil {
		return user
	}

	switch valueType := rawUser.(type) {
	case *userModels.User:
		user = valueType
	}

	return user
}
