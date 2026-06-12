package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateTemplateJawabanCommandValidation(t *testing.T) {
	validUUID1 := "550e8400-e29b-41d4-a716-446655440000"
	validUUID2 := "550e8400-e29b-41d4-a716-446655440001"

	cmd := UpdateTemplateJawabanCommand{
		Uuid:                   validUUID1,
		UuidTemplatePertanyaan: validUUID2,
		Jawaban:                "Tidak",
		IsFreeText:             "0",
		SID:                    "sid-123",
		Resource:               "lpm",
	}

	// Success case
	assert.NoError(t, UpdateTemplateJawabanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, UpdateTemplateJawabanCommandValidation(cmdFailUuid))
}
