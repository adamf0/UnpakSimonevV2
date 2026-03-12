package application

import (
	"context"

	"UnpakSiamida/common/helper"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
)

type CreateTemplateJawabanCommandHandler struct {
	Repo           domaintemplatejawaban.ITemplateJawabanRepository
	RepoPertanyaan domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *CreateTemplateJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTemplateJawabanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidPertanyaan, err := uuid.Parse(cmd.UuidTemplatePertanyaan)
	if err != nil {
		return "", domaintemplatejawaban.InvalidTemplatePertanyaan()
	}

	existingPertanyaan, err := h.RepoPertanyaan.GetByUuid(ctx, uuidPertanyaan)
	if err != nil {
		return "", domaintemplatejawaban.NotFoundTemplatePertanyaan()
	}

	parseFreeText, err := helper.ParseUint(cmd.IsFreeText)
	if err != nil {
		return "", err
	}

	parseNilai, err := helper.ParseUint(cmd.Nilai)
	if err != nil {
		return "", err
	}

	result := domaintemplatejawaban.NewTemplateJawaban(
		existingPertanyaan.ID,
		cmd.Jawaban,
		parseNilai,
		parseFreeText,
		cmd.Resource, //local, mahasiswa, dosen, pegawai
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createTemplateJawaban := result.Value
	if err := h.Repo.Create(ctx, createTemplateJawaban); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
