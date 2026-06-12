package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemplateJawabanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := CreateTemplateJawabanCommand{
		UuidTemplatePertanyaan: validUUID,
		Jawaban:                "Ya",
		IsFreeText:             "0",
		SID:                    "sid-123",
		Resource:               "lpm",
	}

	// Success case
	assert.NoError(t, CreateTemplateJawabanCommandValidation(cmd))

	// Fail on empty UuidTemplatePertanyaan
	cmdFailUuid := cmd
	cmdFailUuid.UuidTemplatePertanyaan = ""
	assert.Error(t, CreateTemplateJawabanCommandValidation(cmdFailUuid))

	// Fail on invalid isFreeText value
	cmdFailFreeText := cmd
	cmdFailFreeText.IsFreeText = "9"
	assert.Error(t, CreateTemplateJawabanCommandValidation(cmdFailFreeText))
}
