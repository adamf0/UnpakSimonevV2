package application

import (
	"context"

	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetBankSoalByUuidQueryHandler struct {
	Repo domainBankSoal.IBankSoalRepository
}

func (h *GetBankSoalByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetBankSoalByUuidQuery,
) (*domainBankSoal.BankSoal, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainBankSoal.NotFound(q.Uuid)
	}

	inBankSoal, err := h.Repo.GetByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainBankSoal.NotFound(q.Uuid)
		}
		return nil, err
	}

	return inBankSoal, nil
}
