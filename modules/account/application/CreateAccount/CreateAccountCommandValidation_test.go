package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountCommandValidation(t *testing.T) {
	cmd := CreateAccountCommand{
		Username: "testuser",
		Password: "password123",
		Level:    "user",
		Name:     "Test User",
	}

	// Success case
	assert.NoError(t, CreateAccountCommandValidation(cmd))

	// Fail on empty Username
	cmdFail := cmd
	cmdFail.Username = ""
	assert.Error(t, CreateAccountCommandValidation(cmdFail))
}
