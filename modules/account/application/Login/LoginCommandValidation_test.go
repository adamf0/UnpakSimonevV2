package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginCommandValidation(t *testing.T) {
	cmd := LoginCommand{
		Username: "user1",
		Password: "password123",
	}

	// Success case
	assert.NoError(t, LoginCommandValidation(cmd))

	// Fail on empty Username
	cmdFailUser := cmd
	cmdFailUser.Username = ""
	assert.Error(t, LoginCommandValidation(cmdFailUser))

	// Fail on empty Password
	cmdFailPass := cmd
	cmdFailPass.Password = ""
	assert.Error(t, LoginCommandValidation(cmdFailPass))
}
