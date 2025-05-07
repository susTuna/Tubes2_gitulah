package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Initialize() error {
	var err error
	pool, err = pgxpool.New(context.Background(), "database")
	return err
}

func Close() {
	pool.Close()
}

func Exec(sql string, args ...any) (pgconn.CommandTag, error) {
	return pool.Exec(context.Background(), sql, args)
}

func Query(sql string, args ...any) (pgx.Rows, error) {
	return pool.Query(context.Background(), sql, args)
}

func QueryRow(sql string, args ...any) pgx.Row {
	return pool.QueryRow(context.Background(), sql, args)
}
