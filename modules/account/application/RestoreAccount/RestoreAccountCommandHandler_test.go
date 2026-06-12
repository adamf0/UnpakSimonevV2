package application

import (
	"context"
	"errors"
	"testing"
	"time"

	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRestoreAccountCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		deletedTime := time.Now()
		existingAccount := &domainaccount.Account{
			UUID:      validUUID,
			DeletedAt: &deletedTime,
		}

		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				assert.Equal(t, validUUID, uid)
				return existingAccount, nil
			},
			UpdateFunc: func(ctx context.Context, account *domainaccount.Account) error {
				assert.Nil(t, account.DeletedAt)
				return nil
			},
		}

		handler := &RestoreAccountCommandHandler{
			Repo: repo,
		}

		cmd := RestoreAccountCommand{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("Failure_InvalidUuid", func(t *testing.T) {
		handler := &RestoreAccountCommandHandler{
			Repo: &mockrepo.MockAccountRepository{},
		}

		cmd := RestoreAccountCommand{
			Uuid: "invalid-uuid",
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

		handler := &RestoreAccountCommandHandler{
			Repo: repo,
		}

		cmd := RestoreAccountCommand{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.NotFound(validUUID.String()).Error(), err.Error())
		assert.Empty(t, res)
	})

	t.Run("Failure_RepoGetError", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return nil, expectedErr
			},
		}

		handler := &RestoreAccountCommandHandler{
			Repo: repo,
		}

		cmd := RestoreAccountCommand{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res)
	})

	t.Run("Failure_RepoUpdateError", func(t *testing.T) {
		deletedTime := time.Now()
		existingAccount := &domainaccount.Account{
			UUID:      validUUID,
			DeletedAt: &deletedTime,
		}
		expectedErr := errors.New("update error")
		repo := &mockrepo.MockAccountRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
				return existingAccount, nil
			},
			UpdateFunc: func(ctx context.Context, account *domainaccount.Account) error {
				return expectedErr
			},
		}

		handler := &RestoreAccountCommandHandler{
			Repo: repo,
		}

		cmd := RestoreAccountCommand{
			Uuid: validUUID.String(),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res)
	})
}
