package application

import (
	"context"
	"errors"
	"testing"

	commondomain "UnpakSiamida/common/domain"
	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetAllKuesionersQueryHandler(t *testing.T) {
	t.Run("Success with Pagination", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domainkuesioner.KuesionerDefault, int64, error) {
				assert.Equal(t, "test-search", search)
				assert.Equal(t, 2, *page)
				assert.Equal(t, 5, *limit)
				assert.False(t, deleted)
				return []domainkuesioner.KuesionerDefault{
					{
						Judul: "Kuesioner A",
					},
				}, 11, nil
			},
		}

		handler := &GetAllKuesionersQueryHandler{
			Repo: repo,
		}

		page := 2
		limit := 5
		q := GetAllKuesionersQuery{
			Search:        "test-search",
			Page:          &page,
			Limit:         &limit,
			Deleted:       false,
			SearchFilters: nil,
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), res.Total)
		assert.Equal(t, 2, res.CurrentPage)
		assert.Equal(t, 3, res.TotalPages) // (11 + 5 - 1) / 5 = 3 pages
		assert.Len(t, res.Data, 1)
		assert.Equal(t, "Kuesioner A", res.Data[0].Judul)
	})

	t.Run("Success default Page and Limit", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domainkuesioner.KuesionerDefault, int64, error) {
				assert.Nil(t, page)
				assert.Nil(t, limit)
				return []domainkuesioner.KuesionerDefault{}, 10, nil
			},
		}

		handler := &GetAllKuesionersQueryHandler{
			Repo: repo,
		}

		q := GetAllKuesionersQuery{}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), res.Total)
		assert.Equal(t, 1, res.CurrentPage)
		assert.Equal(t, 1, res.TotalPages)
	})

	t.Run("Fail Repo GetAll Error", func(t *testing.T) {
		getAllErr := errors.New("get all error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domainkuesioner.KuesionerDefault, int64, error) {
				return nil, 0, getAllErr
			},
		}

		handler := &GetAllKuesionersQueryHandler{
			Repo: repo,
		}

		q := GetAllKuesionersQuery{}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, getAllErr)
	})
}
