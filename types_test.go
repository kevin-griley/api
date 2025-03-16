package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {

	user, err := NewUser("a", "b")

	assert.Nil(t, err)

	fmt.Printf("%+v\n", user)

}
