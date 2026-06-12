package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestoreBankSoalCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := RestoreBankSoalCommand{
		Uuid: validUUID,
	}

	// Success case
	assert.NoError(t, RestoreBankSoalCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFail := cmd
	cmdFail.Uuid = ""
	assert.Error(t, RestoreBankSoalCommandValidation(cmdFail))
}
