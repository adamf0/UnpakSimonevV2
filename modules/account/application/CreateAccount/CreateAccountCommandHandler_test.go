package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/common/helper"
	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountCommandHandler_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockrepo.MockAccountRepository{
			CreateFunc: func(ctx context.Context, account *domainaccount.Account) error {
				return nil
			},
		}

		handler := &CreateAccountCommandHandler{
			Repo: repo,
		}

		cmd := CreateAccountCommand{
			Username: "testuser",
			Password: "password123",
			Level:    "user",
			Name:     "Test User",
			Email:    helper.StrPtr("test@test.com"),
		}

		uid, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, uid)
	})

	t.Run("Failure_RepoError", func(t *testing.T) {
		expectedErr := errors.New("database error")
		repo := &mockrepo.MockAccountRepository{
			CreateFunc: func(ctx context.Context, account *domainaccount.Account) error {
				return expectedErr
			},
		}

		handler := &CreateAccountCommandHandler{
			Repo: repo,
		}

		cmd := CreateAccountCommand{
			Username: "testuser",
			Password: "password123",
			Level:    "user",
			Name:     "Test User",
			Email:    helper.StrPtr("test@test.com"),
		}

		uid, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, uid)
	})
}
