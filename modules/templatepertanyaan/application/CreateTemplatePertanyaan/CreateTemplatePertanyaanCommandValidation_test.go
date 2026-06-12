package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemplatePertanyaanCommandValidation(t *testing.T) {
	validUUID1 := "550e8400-e29b-41d4-a716-446655440000"
	validUUID2 := "550e8400-e29b-41d4-a716-446655440001"

	cmd := CreateTemplatePertanyaanCommand{
		UuidBankSoal: validUUID1,
		Pertanyaan:   "Apakah Anda puas?",
		JenisPilihan: "radio",
		Bobot:        "2",
		UuidKategori: validUUID2,
		Required:     1,
		SID:          "sid-123",
		Resource:     "lpm",
	}

	// Success case
	assert.NoError(t, CreateTemplatePertanyaanCommandValidation(cmd))

	// Fail on empty UuidBankSoal
	cmdFailUuid := cmd
	cmdFailUuid.UuidBankSoal = ""
	assert.Error(t, CreateTemplatePertanyaanCommandValidation(cmdFailUuid))
}
