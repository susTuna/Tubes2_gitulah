package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Initialize() error {
	connstring := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_ADDRESS"),
		os.Getenv("DATABASE_DB"),
	)

	_pool, err := pgxpool.New(context.Background(), connstring)
	pool = _pool

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
