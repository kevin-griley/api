package db

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Psql *sql.DB

func CreateDatabase() (*sql.DB, error) {
	godotenv.Load()
	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Init() error {
	db, err := CreateDatabase()
	if err != nil {
		return err
	}

	Psql = db

	return nil
}
