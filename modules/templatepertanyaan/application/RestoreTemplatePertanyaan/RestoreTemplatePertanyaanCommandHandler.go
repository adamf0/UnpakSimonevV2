package application

import (
	"context"
	"errors"

	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestoreTemplatePertanyaanCommandHandler struct {
	Repo domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *RestoreTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd RestoreTemplatePertanyaanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidTemplatePertanyaan, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingTemplatePertanyaan, err := h.Repo.GetByUuid(ctx, uuidTemplatePertanyaan)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplatepertanyaan.NotFound(cmd.Uuid)
		}
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintemplatepertanyaan.RestoreTemplatePertanyaan(
		existingTemplatePertanyaan,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedTemplatePertanyaan := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTemplatePertanyaan); err != nil {
		return "", err
	}

	return updatedTemplatePertanyaan.UUID.String(), nil
}
