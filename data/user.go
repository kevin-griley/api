package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kevin-griley/api/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserName    string    `json:"user_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	LastRequest time.Time `json:"last_request"`
	IsAdmin     bool      `json:"is_admin"`
	IsVerified  bool      `json:"is_verified"`
	IsDeleted   bool      `json:"is_deleted"`
}

func CreateUser(u *User) (string, error) {
	var id string
	query := `
		INSERT INTO "User" (created_at, updated_at, user_name, email, password, last_request)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := db.Psql.QueryRow(
		query,
		u.CreatedAt,
		u.UpdatedAt,
		u.UserName,
		u.Email,
		u.Password,
		u.LastRequest,
	).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func GetUserByEmail(email string) (*User, error) {
	rows, err := db.Psql.Query(`SELECT * FROM "User" WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", email)
}

func NewUser(email, password string) (*User, error) {

	encpwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &User{
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		UserName:    email,
		Email:       email,
		Password:    string(encpwd),
		LastRequest: time.Now().UTC(),
	}, nil
}

func (usr *User) ValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(usr.Password),
		[]byte(password)) == nil
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.LastRequest,
		&user.IsAdmin,
		&user.IsVerified,
		&user.IsDeleted,
	)
	return user, err
}
