package application

import "UnpakSiamida/common/domain"

type GetAllBankSoalsQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
	Deleted       bool
}
