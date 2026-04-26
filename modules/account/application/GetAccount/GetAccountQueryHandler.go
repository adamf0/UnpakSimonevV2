package application

import (
	"context"

	domainAccount "UnpakSiamida/modules/account/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAccountQueryHandler struct {
	Repo domainAccount.IAccountRepository
}

func (h *GetAccountQueryHandler) Handle(
	ctx context.Context,
	q GetAccountQuery,
) (*domainAccount.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainAccount.NotFound(q.Uuid)
	}

	Account, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainAccount.NotFound(q.Uuid)
		}
		return nil, err
	}

	return Account, nil
}
