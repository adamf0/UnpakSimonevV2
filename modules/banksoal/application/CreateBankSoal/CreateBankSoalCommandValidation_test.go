package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBankSoalCommandValidation(t *testing.T) {
	cmd := CreateBankSoalCommand{
		Judul:    "Soal UTS",
		Semester: "202401",
		SID:      "sid-1",
		Resource: "lpm",
	}

	// Success case
	assert.NoError(t, CreateBankSoalCommandValidation(cmd))

	// Fail on empty Judul
	cmdFailJudul := cmd
	cmdFailJudul.Judul = ""
	assert.Error(t, CreateBankSoalCommandValidation(cmdFailJudul))

	// Fail on invalid Semester format
	cmdFailSemester := cmd
	cmdFailSemester.Semester = "2024"
	assert.Error(t, CreateBankSoalCommandValidation(cmdFailSemester))
}
