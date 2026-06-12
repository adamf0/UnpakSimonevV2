package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyTemplatePertanyaanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := CopyTemplatePertanyaanCommand{
		Uuid:     validUUID,
		SID:      "sid-123",
		Resource: "lpm",
	}

	// Success case
	assert.NoError(t, CopyTemplatePertanyaanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, CopyTemplatePertanyaanCommandValidation(cmdFailUuid))
}
