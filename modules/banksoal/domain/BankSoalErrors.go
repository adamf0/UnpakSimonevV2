package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("BankSoal.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("BankSoal.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("BankSoal.InvalidData", "data is invalid")
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
