package application

import (
	"context"
	"errors"

	"UnpakSiamida/common/helper"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateTemplateJawabanCommandHandler struct {
	Repo           domaintemplatejawaban.ITemplateJawabanRepository
	RepoPertanyaan domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *UpdateTemplateJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateTemplateJawabanCommand,
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
	uuidPertanyaan, err := uuid.Parse(cmd.UuidTemplatePertanyaan)
	if err != nil {
		return "", domaintemplatejawaban.InvalidTemplatePertanyaan()
	}

	parseFreeText, err := helper.ParseUint(cmd.IsFreeText)
	if err != nil {
		return "", err
	}

	parseNilai, err := helper.ParseUint(cmd.Nilai)
	if err != nil {
		return "", err
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

	existingPertanyaan, err := h.RepoPertanyaan.GetByUuid(ctx, uuidPertanyaan)
	if err != nil {
		return "", domaintemplatejawaban.NotFoundTemplatePertanyaan()
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintemplatejawaban.UpdateTemplateJawaban(
		existingTemplateJawaban,
		uuidTemplateJawaban,
		existingPertanyaan.ID,
		cmd.Jawaban,
		parseNilai,
		parseFreeText,
		cmd.Resource, //local, simak, simpeg
		cmd.SID,
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
