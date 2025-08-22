package infrastructure

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"time"
)

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, ErrMissingDSN
	}
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err = conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
