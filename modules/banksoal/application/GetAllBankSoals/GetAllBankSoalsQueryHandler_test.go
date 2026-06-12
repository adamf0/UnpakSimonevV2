package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/common/domain"
	helperCommon "UnpakSiamida/common/helper"
	"UnpakSiamida/modules/banksoal/application/mock"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetAllBankSoalsQueryHandler_Handle(t *testing.T) {
	t.Run("NPM not nil", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &GetAllBankSoalsQueryHandler{Repo: repo}
		cmd := GetAllBankSoalsQuery{
			NPM: helperCommon.StrPtr("12345"),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res.Data)
	})

	t.Run("NIDN not nil", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &GetAllBankSoalsQueryHandler{Repo: repo}
		cmd := GetAllBankSoalsQuery{
			NIDN: helperCommon.StrPtr("12345"),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res.Data)
	})

	t.Run("NIP not nil", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &GetAllBankSoalsQueryHandler{Repo: repo}
		cmd := GetAllBankSoalsQuery{
			NIP: helperCommon.StrPtr("12345"),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res.Data)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetAllFunc: func(
				ctx context.Context,
				search string,
				searchFilters []domain.SearchFilter,
				TargetFakultas string,
				TargetProdi string,
				TargetUnit string,
				TargetStatus string,
				page, limit *int,
				deleted bool,
				active bool,
			) ([]domainBankSoal.BankSoalDefault, int64, error) {
				return nil, 0, errors.New("database error")
			},
		}
		handler := &GetAllBankSoalsQueryHandler{Repo: repo}
		cmd := GetAllBankSoalsQuery{}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res.Data)
	})

	t.Run("success with pagination", func(t *testing.T) {
		limit := 10
		page := 2
		repo := &mock.MockRepository{
			GetAllFunc: func(
				ctx context.Context,
				search string,
				searchFilters []domain.SearchFilter,
				TargetFakultas string,
				TargetProdi string,
				TargetUnit string,
				TargetStatus string,
				pageVal, limitVal *int,
				deleted bool,
				active bool,
			) ([]domainBankSoal.BankSoalDefault, int64, error) {
				assert.Equal(t, limit, *limitVal)
				assert.Equal(t, page, *pageVal)
				return []domainBankSoal.BankSoalDefault{
					{Judul: "Soal 1"},
					{Judul: "Soal 2"},
				}, 25, nil
			},
		}
		handler := &GetAllBankSoalsQueryHandler{Repo: repo}
		cmd := GetAllBankSoalsQuery{
			Limit: &limit,
			Page:  &page,
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, int64(25), res.Total)
		assert.Equal(t, 2, res.CurrentPage)
		assert.Equal(t, 3, res.TotalPages) // ceil(25 / 10) = 3
	})
}
