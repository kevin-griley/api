package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type ContextKey string

const ContextKeyDB ContextKey = "ContextKeyDB"

func CreateDatabase() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL must be set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func Init() (*sql.DB, error) {
	db, err := CreateDatabase()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func WithDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, ContextKeyDB, db)
}

func GetDB(ctx context.Context) (*sql.DB, bool) {
	db, ok := ctx.Value(ContextKeyDB).(*sql.DB)
	return db, ok
}
