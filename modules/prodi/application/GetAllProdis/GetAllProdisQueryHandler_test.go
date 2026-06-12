package application

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/prodi/application/mock"
	domainProdi "UnpakSiamida/modules/prodi/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllProdisQueryHandler_Handle(t *testing.T) {
	prodisList := []domainProdi.ProdiDefault{
		{
			KodeFakultas: "FK-01",
			KodeProdi:    "PR-01",
			NamaProdi:    "Informatika",
		},
		{
			KodeFakultas: "FK-01",
			KodeProdi:    "PR-02",
			NamaProdi:    "Sistem Informasi",
		},
	}

	t.Run("Success case", func(t *testing.T) {
		limit := 10
		page := 1
		repo := &mock.MockProdiRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, pPage, pLimit *int) ([]domainProdi.ProdiDefault, int64, error) {
				assert.Equal(t, "informatika", search)
				assert.Equal(t, &page, pPage)
				assert.Equal(t, &limit, pLimit)
				return prodisList, int64(len(prodisList)), nil
			},
		}

		handler := &GetAllProdisQueryHandler{Repo: repo}
		q := GetAllProdisQuery{
			Search: "informatika",
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
		assert.Equal(t, "PR-01", res.Data[0].KodeProdi)
		assert.Equal(t, "Informatika", res.Data[0].NamaProdi)
	})

	t.Run("Failure case - GetAll error", func(t *testing.T) {
		dbErr := errors.New("database error")
		repo := &mock.MockProdiRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, pPage, pLimit *int) ([]domainProdi.ProdiDefault, int64, error) {
				return nil, 0, dbErr
			},
		}

		handler := &GetAllProdisQueryHandler{Repo: repo}
		q := GetAllProdisQuery{}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
