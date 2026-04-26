package application

import (
	"context"
	"time"

	domainaccount "UnpakSiamida/modules/account/domain"
)

type CreateAccountCommandHandler struct {
	Repo domainaccount.IAccountRepository
}

func (h *CreateAccountCommandHandler) Handle(
	ctx context.Context,
	cmd CreateAccountCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := domainaccount.NewAccount(
		cmd.Username,
		cmd.Password,
		cmd.Level,
		cmd.Name,
		cmd.Email,
		cmd.Fakultas,
		cmd.Prodi,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createAccount := result.Value

	if err := h.Repo.Create(ctx, createAccount); err != nil {
		return "", err
	}

	return createAccount.UUID.String(), nil
}
