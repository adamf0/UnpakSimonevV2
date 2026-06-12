package application

import (
	"context"
	"errors"
	"testing"

	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteAccountCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("Success_SoftDelete", func(t *testing.T) {
		existingAccount := &domainaccount.Account{
			UUID: validUUID,
		}

		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				assert.Equal(t, validUUID, uid)
				return existingAccount, nil
			},
			UpdateFunc: func(ctx context.Context, account *domainaccount.Account) error {
				assert.NotNil(t, account.DeletedAt)
				return nil
			},
		}

		handler := &DeleteAccountCommandHandler{
			Repo: repo,
		}

		cmd := DeleteAccountCommand{
			Uuid: validUUID.String(),
			Mode: "soft_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("Success_HardDelete", func(t *testing.T) {
		existingAccount := &domainaccount.Account{
			UUID: validUUID,
		}

		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return existingAccount, nil
			},
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				assert.Equal(t, validUUID, uid)
				return nil
			},
		}

		handler := &DeleteAccountCommandHandler{
			Repo: repo,
		}

		cmd := DeleteAccountCommand{
			Uuid: validUUID.String(),
			Mode: "hard_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("Failure_InvalidUuid", func(t *testing.T) {
		handler := &DeleteAccountCommandHandler{
			Repo: &mockrepo.MockAccountRepository{},
		}

		cmd := DeleteAccountCommand{
			Uuid: "invalid-uuid",
			Mode: "soft_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.InvalidUuid().Error(), err.Error())
		assert.Empty(t, res)
	})

	t.Run("Failure_NotFound", func(t *testing.T) {
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &DeleteAccountCommandHandler{
			Repo: repo,
		}

		cmd := DeleteAccountCommand{
			Uuid: validUUID.String(),
			Mode: "soft_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.NotFound(validUUID.String()).Error(), err.Error())
		assert.Empty(t, res)
	})

	t.Run("Failure_GetByUuidError", func(t *testing.T) {
		expectedErr := errors.New("db read error")
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return nil, expectedErr
			},
		}

		handler := &DeleteAccountCommandHandler{
			Repo: repo,
		}

		cmd := DeleteAccountCommand{
			Uuid: validUUID.String(),
			Mode: "soft_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res)
	})

	t.Run("Failure_HardDeleteError", func(t *testing.T) {
		existingAccount := &domainaccount.Account{
			UUID: validUUID,
		}
		expectedErr := errors.New("db delete error")
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return existingAccount, nil
			},
			DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
				return expectedErr
			},
		}

		handler := &DeleteAccountCommandHandler{
			Repo: repo,
		}

		cmd := DeleteAccountCommand{
			Uuid: validUUID.String(),
			Mode: "hard_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res)
	})

	t.Run("Failure_SoftDeleteUpdateError", func(t *testing.T) {
		existingAccount := &domainaccount.Account{
			UUID: validUUID,
		}
		expectedErr := errors.New("db update error")
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return existingAccount, nil
			},
			UpdateFunc: func(ctx context.Context, account *domainaccount.Account) error {
				return expectedErr
			},
		}

		handler := &DeleteAccountCommandHandler{
			Repo: repo,
		}

		cmd := DeleteAccountCommand{
			Uuid: validUUID.String(),
			Mode: "soft_delete",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res)
	})
}
