package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/phrasetagg/gofermart/internal/app/db"
	orderErrors "github.com/phrasetagg/gofermart/internal/app/errors/services/order"
	"github.com/phrasetagg/gofermart/internal/app/models"
)

const statusNew = "NEW"
const statusProcessing = "PROCESSING"

//const statusInvalid = "INVALID"
//const statusProcessed = "PROCESSED"

type Order struct {
	DB *db.DB
}

func NewOrderRepository(DB *db.DB) *Order {
	return &Order{DB: DB}
}

func (o *Order) Create(userID int64, orderNumber string) error {
	conn, err := o.DB.GetConn(context.Background())

	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO orders (user_id, number, status, uploaded_at) VALUES ($1,$2,$3,NOW())", userID, orderNumber, statusNew)

	return err
}

func (o *Order) GetOrderByNumber(number string) (models.Order, error) {
	var order models.Order

	conn, err := o.DB.GetConn(context.Background())

	if err != nil {
		return order, err
	}

	err = conn.
		QueryRow(
			context.Background(),
			"SELECT o.id, o.user_id, o.number, o.status, o.uploaded_at, a.value "+
				"FROM orders as o LEFT JOIN accruals a on o.number = a.order_number "+
				"WHERE o.number=$1", number).
		Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &order.Accrual)

	if errors.As(err, &pgx.ErrNoRows) {
		return order, &orderErrors.NotFoundError{}
	}

	return order, err
}

func (o *Order) GetOrdersByUserID(userID int64) ([]models.Order, error) {
	orders := make([]models.Order, 0)

	conn, err := o.DB.GetConn(context.Background())

	if err != nil {
		return orders, err
	}

	rows, err := conn.Query(
		context.Background(),
		"SELECT o.id, o.user_id, o.number, o.status, o.uploaded_at, coalesce(a.value,0) "+
			"FROM orders as o LEFT JOIN accruals a on o.number = a.order_number "+
			"WHERE o.user_id=$1",
		userID)

	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &order.Accrual)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *Order) GetUnprocessedOrders() ([]models.Order, error) {
	orders := make([]models.Order, 0)

	conn, err := o.DB.GetConn(context.Background())

	if err != nil {
		return orders, err
	}

	rows, err := conn.Query(
		context.Background(),
		"SELECT number "+
			"FROM orders "+
			"WHERE status=$1 OR status=$2 ORDER BY uploaded_at",
		statusNew, statusProcessing)

	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Number)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *Order) ProcessOrderAccrual(orderNumber string, status string, accrual float64) error {
	conn, err := o.DB.GetConn(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("TRY TO UPDATE. Onumber: " + orderNumber + " STATUS: " + status + "accrual: ")
	fmt.Println(accrual)

	order, err := o.GetOrderByNumber(orderNumber)
	if err != nil {
		return err
	}

	//_, err = conn.Exec(context.Background(), "UPDATE orders SET status=$1 WHERE number=$2", status, orderNumber)

	_, err = conn.Exec(context.Background(), "INSERT INTO orders (user_id, number, status, uploaded_at) "+
		"VALUES ($1,$2,$3,NOW()) ON CONFLICT (number) DO UPDATE SET status=$1 WHERE number=$2",
		order.UserID,
		orderNumber,
		status,
		status,
		orderNumber,
	)

	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO accruals (user_id, order_number, value, created_at) VALUES ($1,$2,$3,NOW())",
		order.UserID, orderNumber, accrual)

	return err
}
