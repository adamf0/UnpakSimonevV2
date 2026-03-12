package application

import (
	"context"
	"errors"
	"strconv"
	"time"

	"UnpakSiamida/common/helper"
	domainaccount "UnpakSiamida/modules/account/domain"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateKuesionerCommandHandler struct {
	Repo         domainkuesioner.IKuesionerRepository
	RepoBankSoal domainbanksoal.IBankSoalRepository
	RepoAccount  domainaccount.IAccountRepository
}

func (h *CreateKuesionerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateKuesionerCommand,
) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE BANK SOAL
	// -------------------------
	bankUUID, err := uuid.Parse(cmd.UuidBankSoal)
	if err != nil {
		return "", domainkuesioner.InvalidBankSoal()
	}

	bankSoal, err := h.RepoBankSoal.GetByUuid(ctx, bankUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainkuesioner.NotFoundBankSoal()
		}
		return "", err
	}

	// -------------------------
	// RESOLVE ACCOUNT
	// -------------------------
	identifier, err := ResolveAccountIdentifier(cmd)
	if err != nil {
		return "", err
	}

	account, err := h.RepoAccount.Get(ctx, identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainkuesioner.NotFoundBankSoal()
		}
		return "", err
	}

	// -------------------------
	// MAP ACCOUNT -> KUESIONER
	// -------------------------
	k := MapAccountToKuesioner(account, cmd)

	// -------------------------
	// CREATE DOMAIN ENTITY
	// -------------------------
	result := domainkuesioner.NewKuesioner(
		k.NIDN,
		k.NamaDosen,
		k.NIP,
		k.NamaTendik,
		k.NPM,
		k.NamaMahasiswa,
		k.KodeFakultas,
		k.Fakultas,
		k.KodeProdi,
		k.Prodi,
		k.Unit,
		strconv.FormatUint(uint64(bankSoal.ID), 10),
		cmd.Tanggal,
		cmd.Resource,
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	if err := h.Repo.Create(ctx, result.Value); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}

func ResolveAccountIdentifier(
	cmd CreateKuesionerCommand,
) (domainaccount.AccountIdentifier, error) {

	switch cmd.Resource {

	case "simak":

		switch cmd.CodeCtx {

		case domainaccount.CtxDosen:
			return domainaccount.AccountIdentifier{
				NIDN: &cmd.SID,
			}, nil

		case domainaccount.CtxMahasiswa:
			return domainaccount.AccountIdentifier{
				NIM: &cmd.SID,
			}, nil
		}

	case "simpeg":

		return domainaccount.AccountIdentifier{
			NIP: &cmd.SID,
		}, nil

	case "local":
		return domainaccount.AccountIdentifier{}, domainkuesioner.RespondentOnly()
	}

	return domainaccount.AccountIdentifier{}, domainkuesioner.NotFoundResource()
}

type KuesionerAccountData struct {
	NIDN          *string
	NamaDosen     *string
	NIP           *string
	NamaTendik    *string
	NPM           *string
	NamaMahasiswa *string
	KodeFakultas  *string
	Fakultas      *string
	KodeProdi     *string
	Prodi         *string
	Unit          *string
}

func MapAccountToKuesioner(
	acc *domainaccount.Account,
	cmd CreateKuesionerCommand,
) KuesionerAccountData {

	data := KuesionerAccountData{
		Fakultas: helper.StrPtr(helper.StringValue(acc.Fakultas)),
	}

	switch cmd.Resource {

	case "simak":

		if cmd.CodeCtx == domainaccount.CtxDosen {

			data.NIDN = &acc.ID
			data.NamaDosen = acc.Name
			data.Prodi = acc.Prodi
			data.KodeFakultas = acc.RefFakultas
			data.KodeProdi = acc.RefProdi

		} else {

			data.NPM = &acc.ID
			data.NamaMahasiswa = acc.Name
			data.Prodi = acc.Prodi
			data.KodeFakultas = acc.RefFakultas
			data.KodeProdi = acc.RefProdi
		}

	case "simpeg":

		data.NIP = &acc.ID
		data.NamaTendik = acc.Name
		data.Unit = acc.Unit
	}

	return data
}
