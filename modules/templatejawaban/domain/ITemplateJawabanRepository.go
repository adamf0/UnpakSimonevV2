package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITemplateJawabanRepository interface {
	GetByUuid(ctx context.Context, uid uuid.UUID) (*TemplateJawaban, error)
	GetByUUIDs(ctx context.Context, uuids []string) ([]TemplateJawaban, error)
	GetFreeTextByPertanyaan(ctx context.Context, pertanyaanID uint) (*TemplateJawaban, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*TemplateJawabanDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]TemplateJawabanDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *TemplateJawaban) error
	Update(ctx context.Context, aktivitasproker *TemplateJawaban) error
	Delete(ctx context.Context, uid uuid.UUID) error
	CopyByTemplatePertanyaan(
		ctx context.Context,
		tx *gorm.DB,
		sourceTemplatePertanyaanID uint,
		targetTemplatePertanyaanID uint,
	) error
	SetupUuid(ctx context.Context) error

	WithTx(tx any) ITemplateJawabanRepository
	BeginTx(ctx context.Context) (*gorm.DB, error)
}
