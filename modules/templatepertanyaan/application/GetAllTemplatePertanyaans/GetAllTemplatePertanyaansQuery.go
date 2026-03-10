package application

import "UnpakSiamida/common/domain"

type GetAllTemplatePertanyaansQuery struct {
	Search        string
	SearchFilters []domain.SearchFilter
	Page          *int
	Limit         *int
	Deleted       bool
}
