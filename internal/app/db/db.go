package db

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"strings"
)

type DB struct {
	dsn  string
	conn *pgx.Conn
}

func NewDB(dsn string) *DB {
	return &DB{dsn: dsn}
}

func (d *DB) GetConn(ctx context.Context) (*pgx.Conn, error) {
	if d.dsn == "" {
		return nil, errors.New("empty database dsn")
	}

	conn, err := pgx.Connect(ctx, d.dsn)

	if err != nil {
		return nil, err
	}

	d.conn = conn

	return d.conn, nil
}

func (d *DB) Close() error {
	if d.conn == nil {
		return nil
	}

	return d.conn.Close(context.Background())
}

func (d *DB) CreateTables() error {
	conn, err := d.GetConn(context.Background())
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users ("+
		"id SERIAL,"+
		"login text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"password text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"created_at timestamp with time zone,"+
		"CONSTRAINT users_pkey PRIMARY KEY (login)"+
		")")
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "CREATE TYPE order_statuses AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')")
	if err != nil && !strings.Contains(err.Error(), pgerrcode.DuplicateObject) {
		return err
	}

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS orders ("+
		"id SERIAL,"+
		"user_id bigint NOT NULL,"+
		"number text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"status order_statuses,"+
		"uploaded_at timestamp with time zone,"+
		"CONSTRAINT orders_pkey PRIMARY KEY (number)"+
		")")
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS accruals ("+
		"id SERIAL,"+
		"user_id bigint NOT NULL,"+
		"order_number text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"value float NOT NULL,"+
		"created_at timestamp with time zone,"+
		"CONSTRAINT accruals_pkey PRIMARY KEY (id)"+
		")")
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS accruals_withdrawn ("+
		"id SERIAL,"+
		"user_id bigint NOT NULL,"+
		"order_number text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"value float NOT NULL,"+
		"created_at timestamp with time zone,"+
		"CONSTRAINT accruals_withdrawn_pkey PRIMARY KEY (id)"+
		")")
	if err != nil {
		return err
	}

	return d.Close()
}
