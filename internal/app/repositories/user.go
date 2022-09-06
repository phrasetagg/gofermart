package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/phrasetagg/gofermart/internal/app/db"
	userErrors "github.com/phrasetagg/gofermart/internal/app/errors/services/user"
	"github.com/phrasetagg/gofermart/internal/app/models/user"
	"strings"
)

type User struct {
	DB *db.DB
}

func NewUserRepository(DB *db.DB) *User {
	return &User{DB: DB}
}

func (u *User) GetUserByLogin(login string) (*user.User, error) {
	var userModel user.User

	conn, err := u.DB.GetConn(context.Background())
	if err != nil {
		return &userModel, err
	}

	err = conn.
		QueryRow(context.Background(), "SELECT id,login FROM users WHERE login=$1", login).
		Scan(&userModel.ID, &userModel.Login)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return &userModel, &userErrors.NotFoundError{}
	}

	return &userModel, nil
}

func (u *User) GetUserByLoginAndPassword(login string, password string) (*user.User, error) {
	var userModel user.User

	conn, err := u.DB.GetConn(context.Background())

	if err != nil {
		return nil, err
	}

	err = conn.
		QueryRow(context.Background(), "SELECT login FROM users WHERE login=$1 AND password=$2", login, password).
		Scan(&userModel.Login)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &userErrors.NotFoundError{}
	}

	return nil, err
}

func (u *User) Create(login string, password string) error {
	conn, err := u.DB.GetConn(context.Background())

	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO users (login, password) VALUES ($1,$2)", login, password)

	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return &userErrors.AlreadyExistsError{Login: login}
	}

	return err
}
