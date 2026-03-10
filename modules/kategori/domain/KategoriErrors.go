package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Kategori.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("Kategori.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("Kategori.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Kategori.NotFound", fmt.Sprintf("Kategori with identifier %s not found", id))
}

func InvalidHierarchy() domain.Error {
	return domain.NotFoundError("Kategori.InvalidHierarchy", "Kategori invalid hierarchy")
}

func InvalidOwner() domain.Error {
	return domain.NotFoundError("Kategori.InvalidOwner", "only lpm / fakultas / prodi can create bank soal")
}
