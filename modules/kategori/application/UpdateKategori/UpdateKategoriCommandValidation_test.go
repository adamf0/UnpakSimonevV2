package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateKategoriCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := UpdateKategoriCommand{
		Uuid:         validUUID,
		NamaKategori: "Updated",
		SID:          "sid-1",
		Resource:     "resource-1",
	}

	// Success case
	assert.NoError(t, UpdateKategoriCommandValidation(cmd))

	// Fail on empty Uuid
	cmdFailUuid := cmd
	cmdFailUuid.Uuid = ""
	assert.Error(t, UpdateKategoriCommandValidation(cmdFailUuid))

	// Fail on empty NamaKategori
	cmdFailName := cmd
	cmdFailName.NamaKategori = ""
	assert.Error(t, UpdateKategoriCommandValidation(cmdFailName))
}
