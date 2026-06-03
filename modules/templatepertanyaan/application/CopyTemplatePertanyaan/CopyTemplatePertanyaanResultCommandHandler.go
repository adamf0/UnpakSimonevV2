package application

import (
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"context"
)

type CopyTemplatePertanyaanResultCommandHandler struct {
	Repo domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *CopyTemplatePertanyaanResultCommandHandler) Handle(
	ctx context.Context,
	cmd CopyTemplatePertanyaanResultCommand,
) (map[uint]uint, error) {

	return h.Repo.CopyByBankSoal(
		ctx,
		cmd.Tx,
		cmd.SourceBankSoalID,
		cmd.TargetBankSoalID,
		cmd.Resource,
		cmd.Sid,
	)
}
