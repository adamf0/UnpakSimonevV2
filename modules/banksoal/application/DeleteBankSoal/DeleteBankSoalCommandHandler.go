package application

import (
	"context"

	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *DeleteBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteBankSoalCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidBankSoal, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainbanksoal.InvalidUuid()
	}

	// Get existing Aktivitasproker
	existingBankSoal, err := h.Repo.GetByUuid(ctx, uuidBankSoal)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainbanksoal.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if cmd.Mode == "hard_delete" {
		if err := h.Repo.Delete(ctx, uuidBankSoal); err != nil { // FIXED
			return "", err
		}
	} else {
		result := domainbanksoal.DeleteBankSoal(
			existingBankSoal,
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		deleteBankSoal := result.Value

		if err := h.Repo.Update(ctx, deleteBankSoal); err != nil {
			return "", err
		}
	}

	return cmd.Uuid, nil
}
