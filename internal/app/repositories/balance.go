package repositories

import (
	"context"
	"github.com/phrasetagg/gofermart/internal/app/db"
	userModels "github.com/phrasetagg/gofermart/internal/app/models/user"
	"github.com/shopspring/decimal"
)

type Balance struct {
	DB *db.DB
}

func NewBalanceRepository(DB *db.DB) *Balance {
	return &Balance{DB: DB}
}

func (b *Balance) GetUserBalance(userID int64) (*userModels.Balance, error) {
	var balance userModels.Balance

	conn, err := b.DB.GetConn(context.Background())
	if err != nil {
		return &balance, err
	}

	var accruals float64
	err = conn.
		QueryRow(context.Background(), "SELECT coalesce(SUM(value),0) accruals FROM accruals WHERE user_id=$1", userID).
		Scan(&accruals)

	if err != nil {
		return nil, err
	}

	var accrualsWithdrawn float64
	err = conn.
		QueryRow(context.Background(), "SELECT coalesce(SUM(value),0) accruals_withdrawn FROM accruals_withdrawn WHERE user_id=$1", userID).
		Scan(&accrualsWithdrawn)

	if err != nil {
		return nil, err
	}

	accrualsValue := decimal.NewFromFloat(accruals)
	accrualsWithdrawnValue := decimal.NewFromFloat(accrualsWithdrawn)

	balance.Current, _ = accrualsValue.Sub(accrualsWithdrawnValue).Float64()
	balance.Withdrawn, _ = accrualsWithdrawnValue.Float64()

	return &balance, nil
}

func (b *Balance) GetUserWithdrawals(userID int64) ([]userModels.Withdrawal, error) {
	withdrawals := make([]userModels.Withdrawal, 0)

	conn, err := b.DB.GetConn(context.Background())
	if err != nil {
		return withdrawals, err
	}

	rows, err := conn.Query(
		context.Background(),
		"SELECT order_number, coalesce(value,0), created_at "+
			"FROM accruals_withdrawn "+
			"WHERE user_id = $1",
		userID)

	if err != nil {
		return withdrawals, err
	}

	defer rows.Close()

	for rows.Next() {
		var withdrawal userModels.Withdrawal
		err = rows.Scan(&withdrawal.OrderNumber, &withdrawal.Value, &withdrawal.CreatedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (b *Balance) AddWithdraw(userID int64, orderNumber string, withdrawValue float64) error {
	conn, err := b.DB.GetConn(context.Background())
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(),
		"INSERT INTO accruals_withdrawn "+
			"(user_id, order_number, value, created_at) "+
			"VALUES ($1,$2,$3,NOW())",
		userID, orderNumber, withdrawValue)

	return err
}
