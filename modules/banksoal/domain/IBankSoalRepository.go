package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IBankSoalRepository interface {
	CountCopy(ctx context.Context, judul string) (int, error)
	GetByUuid(ctx context.Context, uid uuid.UUID) (*BankSoal, error)
	GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*BankSoalDefault, error)
	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]BankSoalDefault, int64, error)
	Create(ctx context.Context, aktivitasproker *BankSoal) error
	Update(ctx context.Context, aktivitasproker *BankSoal) error
	Delete(ctx context.Context, uid uuid.UUID) error
	SetupUuid(ctx context.Context) error
}
