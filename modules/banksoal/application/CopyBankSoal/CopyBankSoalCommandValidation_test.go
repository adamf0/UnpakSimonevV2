package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyBankSoalCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := CopyBankSoalCommand{
		Uuid:     validUUID,
		SID:      "sid-123",
		Resource: "lpm",
	}

	// Success case
	assert.NoError(t, CopyBankSoalCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, CopyBankSoalCommandValidation(cmdFailUuid))
}
