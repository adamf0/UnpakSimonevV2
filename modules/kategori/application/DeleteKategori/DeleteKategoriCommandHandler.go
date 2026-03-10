package application

import (
	"context"

	domainkategori "UnpakSiamida/modules/kategori/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteKategoriCommandHandler struct {
	Repo domainkategori.IKategoriRepository
}

func (h *DeleteKategoriCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteKategoriCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidKategori, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkategori.InvalidUuid()
	}

	// Get existing Aktivitasproker
	existingKategori, err := h.Repo.GetByUuid(ctx, uuidKategori)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainkategori.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if cmd.Mode == "hard_delete" {
		if err := h.Repo.Delete(ctx, uuidKategori); err != nil { // FIXED
			return "", err
		}
	} else {
		result := domainkategori.DeleteKategori(
			existingKategori,
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		deleteKategori := result.Value

		if err := h.Repo.Update(ctx, deleteKategori); err != nil {
			return "", err
		}
	}

	return cmd.Uuid, nil
}
