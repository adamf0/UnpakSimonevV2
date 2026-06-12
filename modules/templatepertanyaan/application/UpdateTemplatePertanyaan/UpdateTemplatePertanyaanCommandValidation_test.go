package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateTemplatePertanyaanCommandValidation(t *testing.T) {
	validUUID1 := "550e8400-e29b-41d4-a716-446655440000"
	validUUID2 := "550e8400-e29b-41d4-a716-446655440001"
	validUUID3 := "550e8400-e29b-41d4-a716-446655440002"

	cmd := UpdateTemplatePertanyaanCommand{
		Uuid:         validUUID1,
		UuidBankSoal: validUUID2,
		Pertanyaan:   "Apakah Anda puas?",
		JenisPilihan: "radio",
		Bobot:        "2",
		UuidKategori: validUUID3,
		Required:     1,
		SID:          "sid-123",
		Resource:     "lpm",
	}

	// Success case
	assert.NoError(t, UpdateTemplatePertanyaanCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, UpdateTemplatePertanyaanCommandValidation(cmdFailUuid))
}
