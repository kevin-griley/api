package main

import "time"


type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	Name string `json:"name"`
}


type Account struct {
	ID      	string 			`json:"id"`
	CreatedAt 	time.Time 		`json:"createdAt"`
	UpdatedAt 	time.Time 		`json:"updatedAt"`
	Name    	string 			`json:"name"`
	Balance 	int64 			`json:"balance"`
}

func NewAccount (name string) *Account {
	return &Account{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,

	}
}