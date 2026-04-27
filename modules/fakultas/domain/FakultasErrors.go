package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Fakultas.EmptyData", "data is not found")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Fakultas.NotFound", fmt.Sprintf("Fakultas with identifier %s not found", id))
}
