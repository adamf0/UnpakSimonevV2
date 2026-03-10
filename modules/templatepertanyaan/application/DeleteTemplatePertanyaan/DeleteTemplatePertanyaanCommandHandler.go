package application

import (
	"context"

	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteTemplatePertanyaanCommandHandler struct {
	Repo domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *DeleteTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplatePertanyaanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidTemplatePertanyaan, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidUuid()
	}

	// Get existing Aktivitasproker
	existingTemplatePertanyaan, err := h.Repo.GetByUuid(ctx, uuidTemplatePertanyaan)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaintemplatepertanyaan.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if cmd.Mode == "hard_delete" {
		if err := h.Repo.Delete(ctx, uuidTemplatePertanyaan); err != nil { // FIXED
			return "", err
		}
	} else {
		result := domaintemplatepertanyaan.DeleteTemplatePertanyaan(
			existingTemplatePertanyaan,
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		deleteTemplatePertanyaan := result.Value

		if err := h.Repo.Update(ctx, deleteTemplatePertanyaan); err != nil {
			return "", err
		}
	}

	return cmd.Uuid, nil
}
