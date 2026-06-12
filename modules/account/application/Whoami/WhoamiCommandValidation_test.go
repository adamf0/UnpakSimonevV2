package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhoamiCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := WhoamiCommand{
		SID: &validUUID,
	}

	// Success case
	assert.NoError(t, WhoamiCommandValidation(cmd))

	// Fail on empty/nil SID
	cmdFail := cmd
	cmdFail.SID = nil
	assert.Error(t, WhoamiCommandValidation(cmdFail))

	// Fail on empty string SID
	empty := ""
	cmdFail2 := cmd
	cmdFail2.SID = &empty
	assert.Error(t, WhoamiCommandValidation(cmdFail2))
}
