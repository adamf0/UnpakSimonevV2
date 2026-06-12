package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKategoriCommandValidation(t *testing.T) {
	// Success case
	cmd := CreateKategoriCommand{
		NamaKategori: "Test Category",
		SID:          "sid-1",
		Resource:     "resource-1",
	}
	err := CreateKategoriCommandValidation(cmd)
	assert.NoError(t, err)

	// Validation failure: Empty NamaKategori
	cmdFail := cmd
	cmdFail.NamaKategori = ""
	errFail := CreateKategoriCommandValidation(cmdFail)
	assert.Error(t, errFail)
	assert.Contains(t, errFail.Error(), "NamaKategori cannot be blank")

	// Validation failure: Empty SID
	cmdFailSID := cmd
	cmdFailSID.SID = ""
	errFailSID := CreateKategoriCommandValidation(cmdFailSID)
	assert.Error(t, errFailSID)
	assert.Contains(t, errFailSID.Error(), "SID cannot be blank")
}
