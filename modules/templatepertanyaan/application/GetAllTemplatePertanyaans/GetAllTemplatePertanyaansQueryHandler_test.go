package application_test

import (
	"context"
	"errors"
	"testing"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/templatepertanyaan/application/GetAllTemplatePertanyaans"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTemplatePertanyaansQueryHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	page := 1
	limit := 10
	searchFilters := []common.SearchFilter{}
	expectedList := []domain.TemplatePertanyaanDefault{
		{
			ID:         1,
			Pertanyaan: "Pertanyaan 1",
		},
		{
			ID:         2,
			Pertanyaan: "Pertanyaan 2",
		},
	}

	repo.GetAllFunc = func(
		ctx context.Context,
		search string,
		filters []common.SearchFilter,
		pg, lim *int,
		deleted bool,
	) ([]domain.TemplatePertanyaanDefault, int64, error) {
		assert.Equal(t, "test", search)
		assert.Equal(t, &page, pg)
		assert.Equal(t, &limit, lim)
		assert.False(t, deleted)
		return expectedList, 2, nil
	}

	handler := &application.GetAllTemplatePertanyaansQueryHandler{
		Repo: repo,
	}

	q := application.GetAllTemplatePertanyaansQuery{
		Search:        "test",
		SearchFilters: searchFilters,
		Page:          &page,
		Limit:         &limit,
		Deleted:       false,
	}

	res, err := handler.Handle(context.Background(), q)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), res.Total)
	assert.Equal(t, expectedList, res.Data)
	assert.Equal(t, 1, res.CurrentPage)
	assert.Equal(t, 1, res.TotalPages)
}

func TestGetAllTemplatePertanyaansQueryHandler_Fail(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	repo.GetAllFunc = func(
		ctx context.Context,
		search string,
		filters []common.SearchFilter,
		pg, lim *int,
		deleted bool,
	) ([]domain.TemplatePertanyaanDefault, int64, error) {
		return nil, 0, errors.New("db error")
	}

	handler := &application.GetAllTemplatePertanyaansQueryHandler{
		Repo: repo,
	}

	q := application.GetAllTemplatePertanyaansQuery{}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, res.Data)
}
