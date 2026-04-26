package domain

import (
	commonDomain "UnpakSiamida/common/domain"
	"context"

	"github.com/google/uuid"
)

type IAccountRepository interface {
	Auth(ctx context.Context, username string, password string) (*AccountDefault, error)
	Get(ctx context.Context, id AccountIdentifier) (*AccountDefault, error)

	GetByUuid(
		ctx context.Context,
		uid uuid.UUID,
	) (*Account, error)

	GetAll(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		page, limit *int,
		deleted bool,
	) ([]Account, int64, error)

	Create(
		ctx context.Context,
		account *Account,
	) error

	Update(
		ctx context.Context,
		account *Account,
	) error

	Delete(
		ctx context.Context,
		uid uuid.UUID,
	) error

	SetupUuid(
		ctx context.Context,
	) error
}
