package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) (string, error)
	CreateAccount(*Account) (string, error)
	DeleteAccount(string) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(string) (*Account, error)
	GetUserByEmail(string) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgresql://neondb_owner:npg_1jhvYtxEl0On@ep-spring-mode-a6wlk770.us-west-2.aws.neon.tech/neondb?sslmode=require"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {

	if err := s.createAccountTable(); err != nil {
		return err
	}

	if err := s.createUserTable(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS account (
			"id" TEXT NOT NULL DEFAULT gen_random_uuid(),
			"createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"name" TEXT NOT NULL,
			"balance" INT NOT NULL DEFAULT 0,

			CONSTRAINT "account_pkey" PRIMARY KEY ("id")
		)
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS "user" (
			"id" TEXT NOT NULL DEFAULT gen_random_uuid(),
			"createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"email" TEXT NOT NULL,
			"password" TEXT NOT NULL,

			CONSTRAINT "user_pkey" PRIMARY KEY ("id")
		)
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(u *User) (string, error) {
	var id string
	query := `
		INSERT INTO "user" (email, password)
		VALUES ($1, $2)
		RETURNING id
	`

	err := s.db.QueryRow(query, u.Email, u.Password).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) (string, error) {

	var id string
	query := `
		INSERT INTO account (name, balance)
		VALUES ($1, $2)
		RETURNING id
	`

	err := s.db.QueryRow(query, acc.Name, acc.Balance).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *PostgresStore) DeleteAccount(id string) error {
	result, err := s.db.Exec("DELETE FROM account WHERE id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no account found with id: %s", id)
	}

	return nil

}

func (s *PostgresStore) GetAccountByID(id string) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %s not found", id)
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM \"user\" WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", email)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil

}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Name,
		&account.Balance,
	)
	return account, err
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
	)
	return user, err
}
