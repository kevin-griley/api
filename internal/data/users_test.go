package data

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidPassword(t *testing.T) {
	// set up a sample password
	password := "MySecretPassword"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	user := &User{
		HashedPassword: string(hashed),
	}

	// test the valid password case
	if !user.ValidPassword(password) {
		t.Errorf("Expected password to be valid")
	}

	// test with a wrong password
	if user.ValidPassword("wrongpassword") {
		t.Errorf("Expected password to be invalid")
	}
}
