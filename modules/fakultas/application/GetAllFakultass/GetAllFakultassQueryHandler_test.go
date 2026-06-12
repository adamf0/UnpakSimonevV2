package application

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/fakultas/application/mock"
	domainFakultas "UnpakSiamida/modules/fakultas/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllFakultassQueryHandler_Handle(t *testing.T) {
	fakultasList := []domainFakultas.FakultasDefault{
		{
			KodeFakultas: "FK-01",
			NamaFakultas: "Fakultas Teknik",
		},
		{
			KodeFakultas: "FK-02",
			NamaFakultas: "Fakultas Ekonomi",
		},
	}

	t.Run("Success case", func(t *testing.T) {
		limit := 10
		page := 1
		repo := &mock.MockFakultasRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, pPage, pLimit *int) ([]domainFakultas.FakultasDefault, int64, error) {
				assert.Equal(t, "teknik", search)
				assert.Equal(t, &page, pPage)
				assert.Equal(t, &limit, pLimit)
				return fakultasList, int64(len(fakultasList)), nil
			},
		}

		handler := &GetAllFakultassQueryHandler{Repo: repo}
		q := GetAllFakultassQuery{
			Search: "teknik",
			Page:   &page,
			Limit:  &limit,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.Total)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, 1, res.CurrentPage)
		assert.Equal(t, 1, res.TotalPages)
		assert.Equal(t, "FK-01", res.Data[0].KodeFakultas)
		assert.Equal(t, "Fakultas Teknik", res.Data[0].NamaFakultas)
	})

	t.Run("Failure case - GetAll error", func(t *testing.T) {
		dbErr := errors.New("database error")
		repo := &mock.MockFakultasRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, pPage, pLimit *int) ([]domainFakultas.FakultasDefault, int64, error) {
				return nil, 0, dbErr
			},
		}

		handler := &GetAllFakultassQueryHandler{Repo: repo}
		q := GetAllFakultassQuery{}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
