package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteBankSoalCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := DeleteBankSoalCommand{
		Uuid: validUUID,
		Mode: "soft_delete",
	}

	// Success case
	assert.NoError(t, DeleteBankSoalCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, DeleteBankSoalCommandValidation(cmdFailUuid))

	// Fail on empty Mode
	cmdFailMode := cmd
	cmdFailMode.Mode = ""
	errFailMode := DeleteBankSoalCommandValidation(cmdFailMode)
	assert.Error(t, errFailMode)
	assert.Contains(t, errFailMode.Error(), "mode cannot be blank")
}
