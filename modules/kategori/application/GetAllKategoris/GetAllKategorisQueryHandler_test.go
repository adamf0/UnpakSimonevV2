package application

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/kategori/application/mock"
	domainKategori "UnpakSiamida/modules/kategori/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllKategorisQueryHandler_Handle(t *testing.T) {
	kategorisList := []domainKategori.KategoriDefault{
		{
			ID:           1,
			UUID:         uuid.New(),
			NamaKategori: "Category 1",
		},
		{
			ID:           2,
			UUID:         uuid.New(),
			NamaKategori: "Category 2",
		},
	}

	t.Run("Success case", func(t *testing.T) {
		limit := 10
		page := 1
		repo := &mock.MockKategoriRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, pPage, pLimit *int, deleted bool) ([]domainKategori.KategoriDefault, int64, error) {
				assert.Equal(t, "test", search)
				assert.Equal(t, &page, pPage)
				assert.Equal(t, &limit, pLimit)
				assert.False(t, deleted)
				return kategorisList, int64(len(kategorisList)), nil
			},
		}

		handler := &GetAllKategorisQueryHandler{Repo: repo}
		q := GetAllKategorisQuery{
			Search:  "test",
			Page:    &page,
			Limit:   &limit,
			Deleted: false,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.Total)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, 1, res.CurrentPage)
		assert.Equal(t, 1, res.TotalPages)
	})

	t.Run("Failure case - GetAll error", func(t *testing.T) {
		dbErr := errors.New("database error")
		repo := &mock.MockKategoriRepository{
			GetAllFunc: func(ctx context.Context, search string, filters []commonDomain.SearchFilter, pPage, pLimit *int, deleted bool) ([]domainKategori.KategoriDefault, int64, error) {
				return nil, 0, dbErr
			},
		}

		handler := &GetAllKategorisQueryHandler{Repo: repo}
		q := GetAllKategorisQuery{}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}
