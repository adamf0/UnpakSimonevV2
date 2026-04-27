package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"
)

type IProdiRepository interface {
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
	) ([]ProdiDefault, int64, error)
}
