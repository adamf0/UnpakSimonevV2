package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteTimeExtBankSoalCommandValidation(t *testing.T) {
	validUUID1 := "550e8400-e29b-41d4-a716-446655440000"
	validUUID2 := "550e8400-e29b-41d4-a716-446655440001"

	cmd := DeleteTimeExtBankSoalCommand{
		Uuid:         validUUID1,
		UuidBankSoal: validUUID2,
	}

	// Success case
	assert.NoError(t, DeleteTimeExtBankSoalCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, DeleteTimeExtBankSoalCommandValidation(cmdFailUuid))

	// Fail on empty UuidBankSoal
	cmdFailBS := cmd
	cmdFailBS.UuidBankSoal = ""
	assert.Error(t, DeleteTimeExtBankSoalCommandValidation(cmdFailBS))
}
