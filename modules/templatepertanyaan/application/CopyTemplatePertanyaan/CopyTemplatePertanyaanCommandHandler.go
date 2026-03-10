package application

import (
	"context"
	"errors"

	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CopyTemplatePertanyaanCommandHandler struct {
	Repo domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *CopyTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd CopyTemplatePertanyaanCommand,
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
	copyCount, err := h.Repo.CountCopy(ctx, existingTemplatePertanyaan.Pertanyaan)
	if err != nil {
		return "", err
	}

	result := domaintemplatepertanyaan.CopyTemplatePertanyaan(
		existingTemplatePertanyaan,
		copyCount,
		cmd.Resource, //local, mahasiswa, dosen, pegawai
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createTemplatePertanyaan := result.Value
	if err := h.Repo.Create(ctx, createTemplatePertanyaan); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
