package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteTemplateJawabanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := DeleteTemplateJawabanCommand{
		Uuid: validUUID,
		Mode: "soft_delete",
	}

	// Success case
	assert.NoError(t, DeleteTemplateJawabanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, DeleteTemplateJawabanCommandValidation(cmdFailUuid))

	// Fail on empty Mode
	cmdFailMode := cmd
	cmdFailMode.Mode = ""
	errFailMode := DeleteTemplateJawabanCommandValidation(cmdFailMode)
	assert.Error(t, errFailMode)
	assert.Contains(t, errFailMode.Error(), "mode cannot be blank")
}
