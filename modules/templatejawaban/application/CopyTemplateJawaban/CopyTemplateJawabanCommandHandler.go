package application

import (
	"context"

	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
)

type CopyTemplateJawabanCommandHandler struct {
	Repo domaintemplatejawaban.ITemplateJawabanRepository
}

func (h *CopyTemplateJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd CopyTemplateJawabanCommand,
) (string, error) {

	err := h.Repo.CopyByTemplatePertanyaan(
		ctx,
		cmd.Tx,
		cmd.SourceTemplatePertanyaanID,
		cmd.TargetTemplatePertanyaanID,
	)

	if err != nil {
		return "", err
	}

	return "success", nil
}
