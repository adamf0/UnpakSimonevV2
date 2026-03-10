package application

import (
	"context"

	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
)

type SetupUuidBankSoalCommandHandler struct {
	Repo domainBankSoal.IBankSoalRepository
}

func (h *SetupUuidBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidBankSoalCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
