package application

import (
	"context"
	"errors"

	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestoreTemplateJawabanCommandHandler struct {
	Repo domaintemplatejawaban.ITemplateJawabanRepository
}

func (h *RestoreTemplateJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd RestoreTemplateJawabanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidTemplateJawaban, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatejawaban.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingTemplateJawaban, err := h.Repo.GetByUuid(ctx, uuidTemplateJawaban)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplatejawaban.NotFound(cmd.Uuid)
		}
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintemplatejawaban.RestoreTemplateJawaban(
		existingTemplateJawaban,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedTemplateJawaban := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTemplateJawaban); err != nil {
		return "", err
	}

	return updatedTemplateJawaban.UUID.String(), nil
}
