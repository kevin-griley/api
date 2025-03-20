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

func CreateUser(ctx context.Context, u *User) (*User, error) {
	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get db connection")
	}

	data := map[string]any{
		"id":              u.ID,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
		"user_name":       u.UserName,
		"email":           u.Email,
		"hashed_password": u.HashedPassword,
		"last_request":    u.LastRequest,
		"last_login":      u.LastLogin,
	}

	query, values, err := BuildInsertQuery("users", data)
	if err != nil {
		return nil, err
	}

	rows, err := dbConn.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return ScanIntoUser(rows)
	}

	return nil, fmt.Errorf("failed to create user")

}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get db connection")
	}

	data := map[string]any{
		"email": email,
	}

	query, values, err := BuildSelectQuery("users", data)
	if err != nil {
		return nil, err
	}

	rows, err := dbConn.Query(query, values...)
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

	data := map[string]any{
		"id": ID,
	}

	query, values, err := BuildSelectQuery("users", data)
	if err != nil {
		return nil, err
	}

	rows, err := dbConn.Query(query, values...)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return ScanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", ID)
}

func UpdateUser(ctx context.Context, u *User) (*User, error) {

	dbConn, ok := db.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get db connection")
	}

	updateData := make(map[string]any)
	updateData["updated_at"] = u.UpdatedAt

	if u.UserName != "" {
		updateData["user_name"] = u.UserName
	}
	if u.Email != "" {
		updateData["email"] = u.Email
	}
	if u.HashedPassword != "" {
		updateData["hashed_password"] = u.HashedPassword
	}
	if u.IsAdmin {
		updateData["is_admin"] = u.IsAdmin
	}
	if u.IsVerified {
		updateData["is_verified"] = u.IsVerified
	}
	if u.IsDeleted {
		updateData["is_deleted"] = u.IsDeleted
	}
	if !u.LastRequest.IsZero() {
		updateData["last_request"] = u.LastRequest
	}
	if !u.LastLogin.IsZero() {
		updateData["last_login"] = u.LastLogin
	}
	if u.FailedLoginAttempts != 0 {
		updateData["failed_login_attempts"] = u.FailedLoginAttempts
	}

	conditions := map[string]any{
		"id": u.ID,
	}

	query, args, err := BuildUpdateQuery("users", updateData, conditions)
	if err != nil {
		return nil, err
	}

	rows, err := dbConn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return ScanIntoUser(rows)
	}

	return nil, fmt.Errorf("failed to update user")
}

func CreateRequest(Email, Password string) (*User, error) {
	encpwd, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
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
		UserName:       Email,
		Email:          Email,
		HashedPassword: string(encpwd),
		LastRequest:    time.Now().UTC(),
		LastLogin:      time.Now().UTC(),
	}, nil
}

func UpdateRequest(UserName, Password string) (*User, error) {

	user := new(User)

	if Password != "" {
		encpwd, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = string(encpwd)
	}

	if UserName != "" {
		user.UserName = UserName
	}

	return user, nil

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
