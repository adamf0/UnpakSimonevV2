package application

import (
	"context"
	"errors"

	"UnpakSiamida/common/helper"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateTemplatePertanyaanCommandHandler struct {
	Repo         domaintemplatepertanyaan.ITemplatePertanyaanRepository
	RepoKategori domainkategori.IKategoriRepository
	RepoBankSoal domainbanksoal.IBankSoalRepository
}

func (h *UpdateTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateTemplatePertanyaanCommand,
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
	uuidKategori, err := uuid.Parse(cmd.UuidKategori)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidKategori()
	}
	uuidBankSoal, err := uuid.Parse(cmd.UuidBankSoal)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidBankSoal()
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
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintemplatepertanyaan.UpdateTemplatePertanyaan(
		existingTemplatePertanyaan,
		uuidTemplatePertanyaan,
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

	updatedTemplatePertanyaan := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTemplatePertanyaan); err != nil {
		return "", err
	}

	return updatedTemplatePertanyaan.UUID.String(), nil
}
