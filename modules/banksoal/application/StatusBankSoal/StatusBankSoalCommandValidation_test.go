package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusBankSoalCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := StatusBankSoalCommand{
		Uuid:   validUUID,
		Status: "active",
	}

	// Success case
	assert.NoError(t, StatusBankSoalCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, StatusBankSoalCommandValidation(cmdFailUuid))

	// Fail on empty Status
	cmdFailStatus := cmd
	cmdFailStatus.Status = ""
	assert.Error(t, StatusBankSoalCommandValidation(cmdFailStatus))
}
