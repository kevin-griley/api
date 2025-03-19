package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevin-griley/api/internal/db"
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

func CreateUser(ctx context.Context, u *User) (uuid.UUID, error) {
	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return u.ID, fmt.Errorf("failed to get db connection")
	}

	query := `
		INSERT INTO users (id, created_at, updated_at, user_name, email, hashed_password, last_request, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	if err := dbConn.QueryRow(
		query,
		u.ID,
		u.CreatedAt,
		u.UpdatedAt,
		u.UserName,
		u.Email,
		u.HashedPassword,
		u.LastRequest,
		u.LastLogin,
	).Scan(&u.ID); err != nil {
		return u.ID, fmt.Errorf("failed to insert user: %w", err)
	}

	return u.ID, nil
}

func UpdateUser(ctx context.Context, u *User) error {
	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return fmt.Errorf("failed to get db connection")
	}

	u.UpdatedAt = time.Now().UTC()

	const query = `
		UPDATE users SET 
			updated_at = $1,
			user_name = $2,
			email = $3,
			is_admin = $4,
			is_verified = $5,
			is_deleted = $6,
			last_request = $7,
			last_login = $8,
			failed_login_attempts = $9
		WHERE id = $10
	`
	result, err := dbConn.Exec(
		query,
		u.UpdatedAt,
		u.UserName,
		u.Email,
		u.IsAdmin,
		u.IsVerified,
		u.IsDeleted,
		u.LastRequest,
		u.LastLogin,
		u.FailedLoginAttempts,
		u.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %s", u.ID)
	}

	return nil
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get db connection")
	}

	rows, err := dbConn.Query(`SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return ScanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", email)
}

func GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error) {
	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get db connection")
	}

	rows, err := dbConn.Query(`SELECT * FROM users WHERE id = $1`, ID)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return ScanIntoUser(rows)
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

func ScanIntoUser(rows *sql.Rows) (*User, error) {
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
