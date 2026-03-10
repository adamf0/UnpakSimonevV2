package application

import (
	"context"

	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
)

type SetupUuidTemplatePertanyaanCommandHandler struct {
	Repo domainTemplatePertanyaan.ITemplatePertanyaanRepository
}

func (h *SetupUuidTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidTemplatePertanyaanCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
