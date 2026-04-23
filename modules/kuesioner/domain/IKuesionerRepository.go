package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IKuesionerRepository interface {
	GetAllKuesionerResult(
		ctx context.Context,
		JudulBankSoal *string,
		Semester *string,
		Is4Year bool,
	) ([]KuesionerResult, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Kuesioner, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*KuesionerDefault, error)
	GetAllFormFromActiveBankSoal(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]KuesionerDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]KuesionerDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *Kuesioner) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
