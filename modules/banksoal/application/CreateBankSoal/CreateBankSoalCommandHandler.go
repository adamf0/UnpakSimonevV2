package application

import (
	"context"

	"UnpakSiamida/common/helper"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"
)

type CreateBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *CreateBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd CreateBankSoalCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := domainbanksoal.NewBankSoal(
		cmd.Judul,
		helper.StrPtr(cmd.Content),
		helper.StrPtr(cmd.Deskripsi),
		helper.StrPtr(cmd.Semester),
		helper.StrPtr(cmd.TanggalMulai),
		helper.StrPtr(cmd.TanggalAkhir),
		cmd.Resource, //local, simak, simpeg
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createBankSoal := result.Value
	if err := h.Repo.Create(ctx, createBankSoal); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
