package application

import (
	"context"
	"errors"
	"testing"

	commondomain "UnpakSiamida/common/domain"
	mockrepo "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"

	"github.com/stretchr/testify/assert"
)

func intPtr(i int) *int {
	return &i
}

func TestGetAllAccountsQueryHandler_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedAccounts := []domainaccount.Account{
			{ID: 1},
			{ID: 2},
		}

		repo := &mockrepo.MockAccountRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domainaccount.Account, int64, error) {
				return expectedAccounts, 5, nil
			},
		}

		handler := &GetAllAccountsQueryHandler{
			Repo: repo,
		}

		q := GetAllAccountsQuery{
			Search:        "test",
			SearchFilters: []commondomain.SearchFilter{},
			Page:          intPtr(2),
			Limit:         intPtr(2),
			Deleted:       false,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccounts, res.Data)
		assert.Equal(t, int64(5), res.Total)
		assert.Equal(t, 2, res.CurrentPage)
		assert.Equal(t, 3, res.TotalPages) // (5 + 2 - 1) / 2 = 3
	})

	t.Run("Success_NilPageAndLimit", func(t *testing.T) {
		expectedAccounts := []domainaccount.Account{
			{ID: 1},
		}

		repo := &mockrepo.MockAccountRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domainaccount.Account, int64, error) {
				return expectedAccounts, 1, nil
			},
		}

		handler := &GetAllAccountsQueryHandler{
			Repo: repo,
		}

		q := GetAllAccountsQuery{
			Search:  "test",
			Deleted: false,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccounts, res.Data)
		assert.Equal(t, int64(1), res.Total)
		assert.Equal(t, 1, res.CurrentPage)
		assert.Equal(t, 1, res.TotalPages)
	})

	t.Run("Failure_RepoError", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mockrepo.MockAccountRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domainaccount.Account, int64, error) {
				return nil, 0, expectedErr
			},
		}

		handler := &GetAllAccountsQueryHandler{
			Repo: repo,
		}

		q := GetAllAccountsQuery{
			Search:  "test",
			Deleted: false,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, res.Data)
	})
}
