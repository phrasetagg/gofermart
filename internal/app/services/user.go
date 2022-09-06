package services

import (
	userModels "github.com/phrasetagg/gofermart/internal/app/models/user"
	"github.com/phrasetagg/gofermart/internal/app/repositories"
	"sort"
)

type User struct {
	userRepository    *repositories.User
	balanceRepository *repositories.Balance
}

func NewUserService(userRepository *repositories.User, balanceRepository *repositories.Balance) *User {
	return &User{
		userRepository:    userRepository,
		balanceRepository: balanceRepository,
	}
}

func (u User) Login(login string, password string) (*userModels.User, error) {
	return u.userRepository.GetUserByLoginAndPassword(login, password)
}

func (u *User) Register(login string, password string) error {
	return u.userRepository.Create(login, password)
}

func (u *User) GetBalance(userID int64) (*userModels.Balance, error) {
	return u.balanceRepository.GetUserBalance(userID)
}

func (u *User) GetWithdrawals(userID int64) ([]userModels.Withdrawal, error) {
	withdrawals, err := u.balanceRepository.GetUserWithdrawals(userID)
	if err != nil {
		return withdrawals, err
	}

	sort.Sort(byCreatedAt(withdrawals))

	return withdrawals, err
}

func (u *User) RegisterWithdraw(userID int64, orderNumber string, withdrawValue float64) error {
	return u.balanceRepository.AddWithdraw(userID, orderNumber, withdrawValue)
}

type byCreatedAt []userModels.Withdrawal

func (s byCreatedAt) Len() int           { return len(s) }
func (s byCreatedAt) Less(i, j int) bool { return s[j].CreatedAt.After(s[i].CreatedAt) }
func (s byCreatedAt) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
