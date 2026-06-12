package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestoreTemplateJawabanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := RestoreTemplateJawabanCommand{
		Uuid: validUUID,
	}

	// Success case
	assert.NoError(t, RestoreTemplateJawabanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFail := cmd
	cmdFail.Uuid = ""
	assert.Error(t, RestoreTemplateJawabanCommandValidation(cmdFail))
}
