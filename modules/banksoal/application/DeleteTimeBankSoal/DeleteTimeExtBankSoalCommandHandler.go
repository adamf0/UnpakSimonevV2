package application

import (
	"context"

	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteTimeExtBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *DeleteTimeExtBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTimeExtBankSoalCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidParse, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainbanksoal.InvalidUuid()
	}

	uuidBankSoal, err := uuid.Parse(cmd.UuidBankSoal)
	if err != nil {
		return "", domainbanksoal.InvalidUuid()
	}

	existingBankSoal, err := h.Repo.GetByUuid(ctx, uuidBankSoal)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainbanksoal.NotFound(cmd.Uuid)
		}
		return "", err
	}

	if err := h.Repo.DeleteExt(ctx, uuidParse, existingBankSoal.ID); err != nil {
		return "", err
	}

	return cmd.Uuid, nil
}
