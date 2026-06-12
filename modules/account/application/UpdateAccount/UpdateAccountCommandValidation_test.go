package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateAccountCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := UpdateAccountCommand{
		Uuid:     validUUID,
		Username: "testuser",
		Level:    "user",
		Name:     "Test User",
	}

	// Success case
	assert.NoError(t, UpdateAccountCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, UpdateAccountCommandValidation(cmdFailUuid))
}
