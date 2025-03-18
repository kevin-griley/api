package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevin-griley/api/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	UserName            string    `json:"user_name"`
	Email               string    `json:"email"`
	HashedPassword      string    `json:"-"`
	IsAdmin             bool      `json:"-"`
	IsVerified          bool      `json:"-"`
	IsDeleted           bool      `json:"-"`
	LastRequest         time.Time `json:"-"`
	LastLogin           time.Time `json:"-"`
	FailedLoginAttempts int       `json:"-"`
}

func CreateUser(u *User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (id, created_at, updated_at, user_name, email, hashed_password, last_request, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err := db.Psql.QueryRow(
		query,
		u.ID,
		u.CreatedAt,
		u.UpdatedAt,
		u.UserName,
		u.Email,
		u.HashedPassword,
		u.LastRequest,
		u.LastLogin,
	).Scan(&u.ID)
	if err != nil {
		return u.ID, fmt.Errorf("user %s already exists", u.Email)
	}

	return u.ID, nil
}

func GetUserByEmail(email string) (*User, error) {
	rows, err := db.Psql.Query(`SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", email)
}

func GetUserByID(ID string) (*User, error) {
	userId, err := uuid.Parse(ID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Psql.Query(`SELECT * FROM users WHERE id = $1`, userId)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", ID)
}

func NewUser(email, password string) (*User, error) {
	encpwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:             userId,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		UserName:       email,
		Email:          email,
		HashedPassword: string(encpwd),
		LastRequest:    time.Now().UTC(),
		LastLogin:      time.Now().UTC(),
	}, nil
}

func (usr *User) ValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword),
		[]byte(password)) == nil
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	u := new(User)
	err := rows.Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.UserName,
		&u.Email,
		&u.HashedPassword,
		&u.IsAdmin,
		&u.IsVerified,
		&u.IsDeleted,
		&u.LastRequest,
		&u.LastLogin,
		&u.FailedLoginAttempts,
	)
	return u, err
}
