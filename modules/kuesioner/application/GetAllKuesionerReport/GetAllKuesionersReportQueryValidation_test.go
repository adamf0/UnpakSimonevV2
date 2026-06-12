package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllKuesionersReportQueryValidation(t *testing.T) {
	// Success case: Is4Year is true, JudulBankSoal can be nil
	q1 := GetAllKuesionersReportQuery{
		Is4Year: true,
	}
	assert.NoError(t, GetAllKuesionersReportQueryValidation(q1))

	// Success case: Is4Year is false, JudulBankSoal is set
	judul := "Soal Evaluasi"
	q2 := GetAllKuesionersReportQuery{
		Is4Year:       false,
		JudulBankSoal: &judul,
	}
	assert.NoError(t, GetAllKuesionersReportQueryValidation(q2))

	// Failure case: Is4Year is false, JudulBankSoal is nil
	q3 := GetAllKuesionersReportQuery{
		Is4Year:       false,
		JudulBankSoal: nil,
	}
	assert.Error(t, GetAllKuesionersReportQueryValidation(q3))

	// Failure case: Is4Year is false, JudulBankSoal is empty string pointer
	empty := ""
	q4 := GetAllKuesionersReportQuery{
		Is4Year:       false,
		JudulBankSoal: &empty,
	}
	assert.Error(t, GetAllKuesionersReportQueryValidation(q4))
}
