package application

import (
	"context"

	domainAccount "UnpakSiamida/modules/account/domain"
)

type SetupUuidAccountCommandHandler struct {
	Repo domainAccount.IAccountRepository
}

func (h *SetupUuidAccountCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidAccountCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
