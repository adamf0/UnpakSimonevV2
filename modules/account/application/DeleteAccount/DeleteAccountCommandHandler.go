package application

import (
	"context"

	domainaccount "UnpakSiamida/modules/account/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteAccountCommandHandler struct {
	Repo domainaccount.IAccountRepository
}

func (h *DeleteAccountCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteAccountCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidAccount, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainaccount.InvalidUuid()
	}

	// Get existing Aktivitasproker
	existingAccount, err := h.Repo.GetByUuid(ctx, uuidAccount)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainaccount.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// Delete by UUID
	if cmd.Mode == "hard_delete" {
		if err := h.Repo.Delete(ctx, uuidAccount); err != nil { // FIXED
			return "", err
		}
	} else {
		result := domainaccount.DeleteAccount(
			existingAccount,
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		deleteAccount := result.Value

		if err := h.Repo.Update(ctx, deleteAccount); err != nil {
			return "", err
		}
	}

	return cmd.Uuid, nil
}
