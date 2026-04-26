package application

import (
	"context"
	"errors"

	domainaccount "UnpakSiamida/modules/account/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestoreAccountCommandHandler struct {
	Repo domainaccount.IAccountRepository
}

func (h *RestoreAccountCommandHandler) Handle(
	ctx context.Context,
	cmd RestoreAccountCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidAccount, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainaccount.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingAccount, err := h.Repo.GetByUuid(ctx, uuidAccount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainaccount.NotFound(cmd.Uuid)
		}
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainaccount.RestoreAccount(
		existingAccount,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedAccount := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedAccount); err != nil {
		return "", err
	}

	return updatedAccount.UUID.String(), nil
}
