package application

import "UnpakSiamida/common/domain"

type GetAllBankSoalsQuery struct {
	Search         string
	SearchFilter   []domain.SearchFilter
	NPM            *string
	NIDN           *string
	NIP            *string
	TargetFakultas *string
	TargetProdi    *string
	Page           *int
	Limit          *int
	Deleted        bool
}
