package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveKuesionerJawabanCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := SaveKuesionerJawabanCommand{
		UuidKuesioner:  validUUID,
		UuidPertanyaan: validUUID,
		Jawaban:        `["550e8400-e29b-41d4-a716-446655440001"]`,
		SID:            "sid-123",
		Resource:       "lpm",
		CodeCtx:        "dosen",
	}

	// Success case
	assert.NoError(t, SaveKuesionerJawabanCommandValidation(cmd))

	// Fail on empty UuidKuesioner
	cmdFailUuid := cmd
	cmdFailUuid.UuidKuesioner = ""
	assert.Error(t, SaveKuesionerJawabanCommandValidation(cmdFailUuid))
}
