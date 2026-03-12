package application

import (
	"context"

	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteTemplateJawabanCommandHandler struct {
	Repo domaintemplatejawaban.ITemplateJawabanRepository
}

func (h *DeleteTemplateJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplateJawabanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidTemplateJawaban, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatejawaban.InvalidUuid()
	}

	// Get existing Aktivitasproker
	existingTemplateJawaban, err := h.Repo.GetByUuid(ctx, uuidTemplateJawaban)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaintemplatejawaban.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if cmd.Mode == "hard_delete" {
		if err := h.Repo.Delete(ctx, uuidTemplateJawaban); err != nil { // FIXED
			return "", err
		}
	} else {
		result := domaintemplatejawaban.DeleteTemplateJawaban(
			existingTemplateJawaban,
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		deleteTemplateJawaban := result.Value

		if err := h.Repo.Update(ctx, deleteTemplateJawaban); err != nil {
			return "", err
		}
	}

	return cmd.Uuid, nil
}
