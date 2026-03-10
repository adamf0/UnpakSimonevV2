package application

import (
	"context"
	"errors"

	domainkategori "UnpakSiamida/modules/kategori/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestoreKategoriCommandHandler struct {
	Repo domainkategori.IKategoriRepository
}

func (h *RestoreKategoriCommandHandler) Handle(
	ctx context.Context,
	cmd RestoreKategoriCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidKategori, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkategori.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingKategori, err := h.Repo.GetByUuid(ctx, uuidKategori)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainkategori.NotFound(cmd.Uuid)
		}
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainkategori.RestoreKategori(
		existingKategori,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedKategori := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedKategori); err != nil {
		return "", err
	}

	return updatedKategori.UUID.String(), nil
}
