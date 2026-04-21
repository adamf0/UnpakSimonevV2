package application

import (
	"context"
	"errors"

	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StatusBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *StatusBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd StatusBankSoalCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidBankSoal, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainbanksoal.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingBankSoal, err := h.Repo.GetByUuid(ctx, uuidBankSoal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainbanksoal.NotFound(cmd.Uuid)
		}
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainbanksoal.ChangeStatus(
		existingBankSoal,
		cmd.Status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedBankSoal := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedBankSoal); err != nil {
		return "", err
	}

	return updatedBankSoal.UUID.String(), nil
}
