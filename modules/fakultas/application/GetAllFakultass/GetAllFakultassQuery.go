package application

import "UnpakSiamida/common/domain"

type GetAllFakultassQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
}
