package domain

import (
	"UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("BankSoal.EmptyData", "data is not found")
}

func OnlyStudentLecturerStaff() domain.Error {
	return domain.NotFoundError("BankSoal.OnlyStudentLecturerStaff", "only student, lecture and staff can access this feature") //[note] temuan pada security karena memberikan informasi ke external secara spesifik
}

func OnlyAdminFacultyStudyProgram() domain.Error {
	return domain.NotFoundError("BankSoal.OnlyAdminFacultyStudyProgram", "only admin, faculty and study program can access this feature") //[note] temuan pada security karena memberikan informasi ke external secara spesifik
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("BankSoal.InvalidUuid", "uuid is invalid")
}

func InvalidData(message ...*string) domain.Error {
	var msg *string

	if len(message) > 0 {
		msg = message[0]
	}

	if helper.NullableString(msg) == "" {
		msg = helper.StrPtr("data is invalid")
	}

	return domain.NotFoundError(
		"BankSoal.InvalidData",
		helper.NullableString(msg),
	)
}

func InvalidStatus() domain.Error {
	return domain.NotFoundError("BankSoal.InvalidStatus", "status is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("BankSoal.NotFound", fmt.Sprintf("BankSoal with identifier %s not found", id))
}

func InvalidDate(target string) domain.Error {
	return domain.NotFoundError("BankSoal.InvalidDate", fmt.Sprintf("%s period have wrong date format", target))
}

func InvalidDateRange() domain.Error {
	return domain.NotFoundError("BankSoal.InvalidDateRange", "tanggal akhir must not be earlier than tanggal awal")
}

func InvalidOwner() domain.Error {
	return domain.NotFoundError("BankSoal.InvalidOwner", "only lpm / fakultas / prodi can create bank soal")
}
