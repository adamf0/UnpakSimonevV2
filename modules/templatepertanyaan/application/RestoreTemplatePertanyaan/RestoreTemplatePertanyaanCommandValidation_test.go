package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestoreTemplatePertanyaanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := RestoreTemplatePertanyaanCommand{
		Uuid: validUUID,
	}

	// Success case
	assert.NoError(t, RestoreTemplatePertanyaanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFail := cmd
	cmdFail.Uuid = ""
	assert.Error(t, RestoreTemplatePertanyaanCommandValidation(cmdFail))
}
