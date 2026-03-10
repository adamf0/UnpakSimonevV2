package application

import (
	"context"
	"errors"

	"UnpakSiamida/common/helper"
	domainaccount "UnpakSiamida/modules/account/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WhoamiCommandHandler struct {
	Repo domainaccount.IAccountRepository
}

func (h *WhoamiCommandHandler) Handle(
	ctx context.Context,
	cmd WhoamiCommand,
) (*domainaccount.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if cmd.SID != nil {
		_, err := uuid.Parse(helper.StringValue(cmd.SID))
		if err != nil {
			return nil, domainaccount.NotFound(helper.StringValue(cmd.SID))
		}
	}

	account := domainaccount.AccountIdentifier{
		UserID: cmd.SID,
		NIDN:   cmd.NIDN,
		NIP:    cmd.NIP,
		NIM:    cmd.NIM,
	}
	user, err := h.Repo.Get(ctx, account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainaccount.InvalidCredential()
		}
		return nil, err
	}

	return user, nil
}
