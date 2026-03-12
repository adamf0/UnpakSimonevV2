package application

import (
	"context"
	"errors"

	helper "UnpakSiamida/common/helper"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *UpdateBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateBankSoalCommand,
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
	result := domainbanksoal.UpdateBankSoal(
		existingBankSoal,
		uuidBankSoal,
		cmd.Judul,
		helper.StrPtr(cmd.Content),
		helper.StrPtr(cmd.Deskripsi),
		helper.StrPtr(cmd.Semester),
		helper.StrPtr(cmd.TanggalMulai),
		helper.StrPtr(cmd.TanggalAkhir),
		cmd.Resource, //local, simak, simpeg
		cmd.SID,
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
