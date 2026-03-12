package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IKuesionerRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Kuesioner, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*KuesionerDefault, error)
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
