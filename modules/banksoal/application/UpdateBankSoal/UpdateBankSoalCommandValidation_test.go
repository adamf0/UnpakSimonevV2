package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBankSoalCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := UpdateBankSoalCommand{
		Uuid:     validUUID,
		Judul:    "Soal UAS",
		Semester: "202402",
		SID:      "sid-123",
		Resource: "lpm",
	}

	// Success case
	assert.NoError(t, UpdateBankSoalCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, UpdateBankSoalCommandValidation(cmdFailUuid))
}
