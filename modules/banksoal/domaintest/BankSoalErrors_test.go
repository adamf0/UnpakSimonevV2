package domaintest

import (
	"UnpakSiamida/modules/banksoal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBankSoalErrors(t *testing.T) {
	assert.Equal(t, "BankSoal.EmptyData", domain.EmptyData().Code)
	assert.Equal(t, "BankSoal.OnlyStudentLecturerStaff", domain.OnlyStudentLecturerStaff().Code)
	assert.Equal(t, "BankSoal.OnlyAdminFacultyStudyProgram", domain.OnlyAdminFacultyStudyProgram().Code)
	assert.Equal(t, "BankSoal.InvalidUuid", domain.InvalidUuid().Code)
	assert.Equal(t, "BankSoal.InvalidData", domain.InvalidData().Code)
	assert.Equal(t, "BankSoal.InvalidStatus", domain.InvalidStatus().Code)
	assert.Equal(t, "BankSoal.NotFound", domain.NotFound("123").Code)
	assert.Equal(t, "BankSoal.InvalidDate", domain.InvalidDate("start").Code)
	assert.Equal(t, "BankSoal.InvalidDateRange", domain.InvalidDateRange().Code)
	assert.Equal(t, "BankSoal.InvalidOwner", domain.InvalidOwner().Code)
}
