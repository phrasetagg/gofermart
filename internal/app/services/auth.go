package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const AuthTokenName = "AuthToken"
const secretKey = "my perfect gophermart"

type Auth struct{}

func NewAuthService() *Auth {
	return &Auth{}
}

func (a *Auth) ValidateAuthToken(authToken string) bool {
	var (
		data []byte // декодированное сообщение с подписью
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(authToken)
	if err != nil {
		return false
	}

	if len(data) < len(data)-32 {
		return false
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data[:len(data)-32])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[len(data)-32:]) {
		return true
	}

	return false
}

func (a *Auth) GetUserLoginFromAuthToken(authToken string) string {
	data, _ := hex.DecodeString(authToken)

	return string(data[:len(data)-32])
}

func (a *Auth) GenerateAuthToken(login string) string {
	loginBytes := []byte(login)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(loginBytes)
	sign := h.Sum(nil)

	authToken := append(loginBytes, sign...)

	return hex.EncodeToString(authToken)
}
