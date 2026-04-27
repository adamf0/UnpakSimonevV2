package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Prodi.EmptyData", "data is not found")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Prodi.NotFound", fmt.Sprintf("Prodi with identifier %s not found", id))
}
