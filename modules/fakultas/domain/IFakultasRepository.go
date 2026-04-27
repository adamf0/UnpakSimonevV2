package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"
)

type IFakultasRepository interface {
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]FakultasDefault, int64, error)
}
