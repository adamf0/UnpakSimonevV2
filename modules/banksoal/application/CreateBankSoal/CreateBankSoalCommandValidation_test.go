package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBankSoalCommandValidation(t *testing.T) {
	cmd := CreateBankSoalCommand{
		Judul:      "Soal UTS",
		Semester:   "202401",
		Peruntukan: "mahasiswa",
		SID:        "sid-1",
		Resource:   "lpm",
	}

	// Success case
	assert.NoError(t, CreateBankSoalCommandValidation(cmd))

	// Fail on empty Judul
	cmdFailJudul := cmd
	cmdFailJudul.Judul = ""
	assert.Error(t, CreateBankSoalCommandValidation(cmdFailJudul))

	// Fail on empty Semester
	cmdFailSemester := cmd
	cmdFailSemester.Semester = ""
	assert.Error(t, CreateBankSoalCommandValidation(cmdFailSemester))

	// Fail on empty Peruntukan
	cmdFailPeruntukan := cmd
	cmdFailPeruntukan.Peruntukan = ""
	assert.Error(t, CreateBankSoalCommandValidation(cmdFailPeruntukan))
}
