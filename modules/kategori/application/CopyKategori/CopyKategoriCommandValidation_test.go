package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyKategoriCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := CopyKategoriCommand{
		Uuid:     validUUID,
		SID:      "sid-1",
		Resource: "resource-1",
	}

	// Success
	assert.NoError(t, CopyKategoriCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFail := cmd
	cmdFail.Uuid = ""
	assert.Error(t, CopyKategoriCommandValidation(cmdFail))

	// Fail on invalid UUID format
	cmdFail.Uuid = "not-a-uuid"
	assert.Error(t, CopyKategoriCommandValidation(cmdFail))
}
