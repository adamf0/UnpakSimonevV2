package application

import (
	"context"
	"errors"

	"UnpakSiamida/common/helper"
	domainaccount "UnpakSiamida/modules/account/domain"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTemplatePertanyaanCommandHandler struct {
	Repo         domaintemplatepertanyaan.ITemplatePertanyaanRepository
	RepoKategori domainkategori.IKategoriRepository
	RepoBankSoal domainbanksoal.IBankSoalRepository
	RepoAccount  domainaccount.IAccountRepository
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

	identity, err := cmd.AccountIdentifierCreate()
	if err != nil {
		return "", err
	}

	var refFakultas *string
	var refProdi *string
	if h.RepoAccount != nil {
		account, err := h.RepoAccount.Get(ctx, identity, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", domaintemplatepertanyaan.InvalidIdentity()
			}
			return "", err
		}
		if account != nil {
			refFakultas = account.RefFakultas
			refProdi = account.RefProdi
		}
	}

	result := domaintemplatepertanyaan.NewTemplatePertanyaan(
		existingBankSoal.ID,
		cmd.Pertanyaan,
		cmd.JenisPilihan,
		parseBobot,
		&existingKategori.ID,
		cmd.Required,
		refFakultas,
		refProdi,
		cmd.Resource, //local, simak, simpeg
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

func (cmd *CreateTemplatePertanyaanCommand) AccountIdentifierCreate() (domainaccount.AccountIdentifier, error) {
	switch cmd.Resource {
	case "simak":
		return domainaccount.AccountIdentifier{
			NIDN: &cmd.SID,
		}, nil

	case "simpeg":
		return domainaccount.AccountIdentifier{
			NIP: &cmd.SID,
		}, nil

	case "local":
		return domainaccount.AccountIdentifier{
			UserID: &cmd.SID,
		}, nil

	default:
		return domainaccount.AccountIdentifier{}, domaintemplatepertanyaan.InvalidIdentity()
	}
}
