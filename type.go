package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	Name string `json:"name"`
}

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Account struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
}

func NewAccount(name string) *Account {
	return &Account{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
}

func (usr *User) ValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(usr.Password),
		[]byte(password)) == nil
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
