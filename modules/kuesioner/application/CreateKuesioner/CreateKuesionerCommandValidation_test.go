package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKuesionerCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := CreateKuesionerCommand{
		UuidBankSoal: validUUID,
		Tanggal:      "2024-01-01 10:20:30",
		SID:          "sid-123",
		Resource:     "lpm",
		CodeCtx:      "dosen",
	}

	// Success case
	assert.NoError(t, CreateKuesionerCommandValidation(cmd))

	// Fail on empty UuidBankSoal
	cmdFailUuid := cmd
	cmdFailUuid.UuidBankSoal = ""
	assert.Error(t, CreateKuesionerCommandValidation(cmdFailUuid))
}
