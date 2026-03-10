package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IKategoriRepository interface {
	CountCopy(ctx context.Context, judul string) (int, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*Kategori, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*KategoriDefault, error)
	GetChildren(ctx context.Context, parentID int) ([]Kategori, error)
	RebuildFullText(ctx context.Context) error
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]KategoriDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *Kategori) error
	Update(ctx context.Context, aktivitasproker *Kategori) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error

	WithTx(ctx context.Context, fn func(repo IKategoriRepository) error) error
}
