package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kevin-griley/api/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

func CreateUser(u *User) (string, error) {
	var id string
	query := `
		INSERT INTO "user" (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	err := db.Psql.QueryRow(query, u.Email, u.Password).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func GetUserByEmail(email string) (*User, error) {
	rows, err := db.Psql.Query("SELECT * FROM \"user\" WHERE email = $1", email)
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     email,
		Password:  string(encpwd),
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
		&user.Email,
		&user.Password,
	)
	return user, err
}
