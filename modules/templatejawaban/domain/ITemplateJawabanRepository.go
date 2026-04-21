package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
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
	SetupUuid(ctx context.Context) error
}
