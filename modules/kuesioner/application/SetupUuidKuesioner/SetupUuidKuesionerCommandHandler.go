package application

import (
	"context"

	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
)

type SetupUuidKuesionerCommandHandler struct {
	Repo domainKuesioner.IKuesionerRepository
}

func (h *SetupUuidKuesionerCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidKuesionerCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
