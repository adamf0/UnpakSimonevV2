package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusTemplatePertanyaanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := StatusTemplatePertanyaanCommand{
		Uuid:   validUUID,
		Status: "active",
	}

	// Success case
	assert.NoError(t, StatusTemplatePertanyaanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, StatusTemplatePertanyaanCommandValidation(cmdFailUuid))
}
