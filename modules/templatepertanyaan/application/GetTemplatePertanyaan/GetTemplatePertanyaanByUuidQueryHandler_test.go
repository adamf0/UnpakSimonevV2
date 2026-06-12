package application_test

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetTemplatePertanyaanByUuidQueryHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	expectedTp := &domain.TemplatePertanyaan{
		ID:         1,
		UUID:       uid,
		Pertanyaan: "Test Pertanyaan",
	}

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		assert.Equal(t, uid, id)
		return expectedTp, nil
	}

	handler := &application.GetTemplatePertanyaanByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.NoError(t, err)
	assert.Equal(t, expectedTp, res)
}

func TestGetTemplatePertanyaanByUuidQueryHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.GetTemplatePertanyaanByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanByUuidQuery{
		Uuid: "invalid-uuid",
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, res)
}

func TestGetTemplatePertanyaanByUuidQueryHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.GetTemplatePertanyaanByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, res)
}

func TestGetTemplatePertanyaanByUuidQueryHandler_DbError(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaan, error) {
		return nil, errors.New("db error")
	}

	handler := &application.GetTemplatePertanyaanByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Nil(t, res)
}
