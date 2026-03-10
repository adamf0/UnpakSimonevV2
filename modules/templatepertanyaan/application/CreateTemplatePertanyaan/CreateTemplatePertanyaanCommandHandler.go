package application

import (
	"context"

	"UnpakSiamida/common/helper"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
)

type CreateTemplatePertanyaanCommandHandler struct {
	Repo         domaintemplatepertanyaan.ITemplatePertanyaanRepository
	RepoKategori domainkategori.IKategoriRepository
	RepoBankSoal domainbanksoal.IBankSoalRepository
}

func (h *CreateTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTemplatePertanyaanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidKategori, err := uuid.Parse(cmd.UuidKategori)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidKategori()
	}
	uuidBankSoal, err := uuid.Parse(cmd.UuidBankSoal)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidBankSoal()
	}

	existingKategori, err := h.RepoKategori.GetByUuid(ctx, uuidKategori)
	if err != nil {
		return "", domaintemplatepertanyaan.NotFoundKategori()
	}

	existingBankSoal, err := h.RepoBankSoal.GetByUuid(ctx, uuidBankSoal)
	if err != nil {
		return "", domaintemplatepertanyaan.NotFoundBankSoal()
	}

	parseBobot, err := helper.ParseUint(cmd.Bobot)
	if err != nil {
		return "", err
	}

	result := domaintemplatepertanyaan.NewTemplatePertanyaan(
		existingBankSoal.ID,
		cmd.Pertanyaan,
		cmd.JenisPilihan,
		parseBobot,
		&existingKategori.ID,
		cmd.Required,
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
