package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteKategoriCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := DeleteKategoriCommand{
		Uuid: validUUID,
		Mode: "soft_delete",
	}

	// Success case
	assert.NoError(t, DeleteKategoriCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, DeleteKategoriCommandValidation(cmdFailUuid))

	// Fail on empty Mode
	cmdFailMode := cmd
	cmdFailMode.Mode = ""
	errFailMode := DeleteKategoriCommandValidation(cmdFailMode)
	assert.Error(t, errFailMode)
	assert.Contains(t, errFailMode.Error(), "mode cannot be blank")
}
