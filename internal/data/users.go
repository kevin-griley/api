package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *userStoreImpl) CreateUser(u *User) (*User, error) {

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

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("failed to create user")

}

func (s *userStoreImpl) GetUserByEmail(email string) (*User, error) {

	data := map[string]any{
		"email": email,
	}

	query, values, err := BuildSelectQuery("users", data)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", email)
}

func (s *userStoreImpl) GetUserByID(ID uuid.UUID) (*User, error) {

	data := map[string]any{
		"id": ID,
	}

	query, values, err := BuildSelectQuery("users", data)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoUser(rows)
	}
	return nil, fmt.Errorf("user %s not found", ID)
}

func (s *userStoreImpl) UpdateUser(u *User) (*User, error) {

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

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("failed to update user")
}

func (s *userStoreImpl) CreateRequest(Email, Password string) (*User, error) {
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

func (s *userStoreImpl) UpdateRequest(UserName, Password string) (*User, error) {

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

type userStoreImpl struct {
	db *sql.DB
}

var NewUserStore = func(db *sql.DB) UserStore {
	return &userStoreImpl{
		db: db,
	}
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(ID uuid.UUID) (*User, error)

	CreateUser(user *User) (*User, error)
	CreateRequest(email, password string) (*User, error)

	UpdateUser(user *User) (*User, error)
	UpdateRequest(userName, password string) (*User, error)
}

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
