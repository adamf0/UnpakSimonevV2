package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type ITemplatePertanyaanRepository interface {
	CountCopy(ctx context.Context, judul string) (int, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*TemplatePertanyaan, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*TemplatePertanyaanDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]TemplatePertanyaanDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *TemplatePertanyaan) error
	Update(ctx context.Context, aktivitasproker *TemplatePertanyaan) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
