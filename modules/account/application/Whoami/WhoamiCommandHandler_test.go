package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/common/helper"
	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestWhoamiCommandHandler_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedUser := &domainaccount.AccountDefault{
			ID:   "user-123",
			Name: helper.StrPtr("Test User"),
		}

		repo := &mockrepo.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				assert.Equal(t, "user-123", *id.UserID)
				assert.Equal(t, "nim-123", *id.NIM)
				return expectedUser, nil
			},
		}

		handler := &WhoamiCommandHandler{
			Repo: repo,
		}

		cmd := WhoamiCommand{
			SID: helper.StrPtr("user-123"),
			NIM: helper.StrPtr("nim-123"),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, res)
	})

	t.Run("Failure_NotFound_InvalidCredential", func(t *testing.T) {
		repo := &mockrepo.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		handler := &WhoamiCommandHandler{
			Repo: repo,
		}

		cmd := WhoamiCommand{
			SID: helper.StrPtr("user-123"),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainaccount.InvalidCredential().Error(), err.Error())
		assert.Nil(t, res)
	})

	t.Run("Failure_RepoError", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mockrepo.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				return nil, expectedErr
			},
		}

		handler := &WhoamiCommandHandler{
			Repo: repo,
		}

		cmd := WhoamiCommand{
			SID: helper.StrPtr("user-123"),
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
	})
}
