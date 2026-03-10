package application

import (
	"context"

	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetBankSoalDefaultByUuidQueryHandler struct {
	Repo domainBankSoal.IBankSoalRepository
}

func (h *GetBankSoalDefaultByUuidQueryHandler) Handle(
	ctx context.Context,
	q GetBankSoalDefaultByUuidQuery,
) (*domainBankSoal.BankSoalDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsed, err := uuid.Parse(q.Uuid)
	if err != nil {
		return nil, domainBankSoal.NotFound(q.Uuid)
	}

	BankSoal, err := h.Repo.GetDefaultByUuid(ctx, parsed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainBankSoal.NotFound(q.Uuid)
		}
		return nil, err
	}

	return BankSoal, nil
}
