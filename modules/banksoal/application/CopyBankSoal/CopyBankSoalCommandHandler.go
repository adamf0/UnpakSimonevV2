package application

import (
	"context"
	"errors"

	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CopyBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *CopyBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd CopyBankSoalCommand,
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
	copyCount, err := h.Repo.CountCopy(ctx, existingBankSoal.Judul)
	if err != nil {
		return "", err
	}

	result := domainbanksoal.CopyBankSoal(
		existingBankSoal,
		copyCount,
		cmd.Resource, //local, mahasiswa, dosen, pegawai
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createBankSoal := result.Value
	if err := h.Repo.Create(ctx, createBankSoal); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
