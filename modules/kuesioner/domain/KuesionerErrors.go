package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Kuesioner.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("Kuesioner.InvalidUuid", "uuid is invalid")
}

func InvalidBankSoal() domain.Error {
	return domain.NotFoundError("Kuesioner.InvalidBankSoal", "bank soal is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Kuesioner.NotFound", fmt.Sprintf("kuesioner with identifier %s not found", id))
}

func NotFoundBankSoal() domain.Error {
	return domain.NotFoundError("Kuesioner.NotFoundBankSoal", "bank soal is not found")
}

func NotFoundResource() domain.Error {
	return domain.NotFoundError("Kuesioner.NotFoundResource", "resource is not found")
}

func NoInfoAccount() domain.Error {
	return domain.NotFoundError("Kuesioner.NoInfoAccount", "no information for account")
}

func RespondentOnly() domain.Error {
	return domain.NotFoundError("Kuesioner.RespondentOnly", "only respondents may fill in")
}

func InvalidDate(target string) domain.Error {
	return domain.NotFoundError("Kuesioner.InvalidDate", fmt.Sprintf("%s period have wrong date format", target))
}

func InvalidDateRange() domain.Error {
	return domain.NotFoundError("Kuesioner.InvalidDateRange", "tanggal akhir must not be earlier than tanggal awal")
}

func Expired() domain.Error {
	return domain.NotFoundError("Kuesioner.Expired", "kuesioner is expired")
}
