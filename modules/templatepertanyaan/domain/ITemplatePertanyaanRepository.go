package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITemplatePertanyaanRepository interface {
	CountCopy(ctx context.Context, judul string) (int, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*TemplatePertanyaan, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*TemplatePertanyaanDefault, error)
	GetDefaultWithAnswareByUuid(ctx context.Context, uid uuid.UUID) (*TemplatePertanyaanWithAnswareDefault, error)
	GetDefaultWithAnswareByBankSoal(ctx context.Context, id_banksoal uint) ([]TemplatePertanyaanWithAnswareDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]TemplatePertanyaanDefault, int64, error)
	CopyByBankSoal(
		ctx context.Context,
		tx *gorm.DB,
		sourceBankSoalID uint,
		targetBankSoalID uint,
		resource string,
		sid string,
	) (map[uint]uint, error)
	Create(ctx context.Context, aktivitasproker *TemplatePertanyaan) error
	Update(ctx context.Context, aktivitasproker *TemplatePertanyaan) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error

	WithTx(tx any) ITemplatePertanyaanRepository
	BeginTx(ctx context.Context) (*gorm.DB, error)
}
