package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.InvalidUuid", "uuid is invalid")
}

func InvalidKategori() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.InvalidKategori", "kategori is invalid")
}

func InvalidBankSoal() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.InvalidBankSoal", "kategori is invalid")
}

func NotFoundKategori() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.NotFoundKategori", "kategori is not found")
}
func NotFoundBankSoal() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.NotFoundBankSoal", "kategori is not found")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.NotFound", fmt.Sprintf("TemplatePertanyaan with identifier %s not found", id))
}

func InvalidOwner() domain.Error {
	return domain.NotFoundError("TemplatePertanyaan.InvalidOwner", "only lpm / fakultas / prodi can create bank soal")
}
