package application

import (
	"context"

	domainTemplateJawaban "UnpakSiamida/modules/templatejawaban/domain"
)

type SetupUuidTemplateJawabanCommandHandler struct {
	Repo domainTemplateJawaban.ITemplateJawabanRepository
}

func (h *SetupUuidTemplateJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidTemplateJawabanCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
