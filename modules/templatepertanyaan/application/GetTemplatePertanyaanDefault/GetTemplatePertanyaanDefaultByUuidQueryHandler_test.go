package application_test

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaanDefault"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetTemplatePertanyaanDefaultByUuidQueryHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	expectedTp := &domain.TemplatePertanyaanDefault{
		ID:         1,
		UUID:       uid,
		Pertanyaan: "Test Default Pertanyaan",
	}

	repo.GetDefaultByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaanDefault, error) {
		assert.Equal(t, uid, id)
		return expectedTp, nil
	}

	handler := &application.GetTemplatePertanyaanDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanDefaultByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.NoError(t, err)
	assert.Equal(t, expectedTp, res)
}

func TestGetTemplatePertanyaanDefaultByUuidQueryHandler_InvalidUuid(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	handler := &application.GetTemplatePertanyaanDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanDefaultByUuidQuery{
		Uuid: "invalid-uuid",
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, res)
}

func TestGetTemplatePertanyaanDefaultByUuidQueryHandler_NotFound(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetDefaultByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaanDefault, error) {
		return nil, gorm.ErrRecordNotFound
	}

	handler := &application.GetTemplatePertanyaanDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanDefaultByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, res)
}

func TestGetTemplatePertanyaanDefaultByUuidQueryHandler_DbError(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}
	uid := uuid.New()

	repo.GetDefaultByUuidFunc = func(ctx context.Context, id uuid.UUID) (*domain.TemplatePertanyaanDefault, error) {
		return nil, errors.New("db error")
	}

	handler := &application.GetTemplatePertanyaanDefaultByUuidQueryHandler{
		Repo: repo,
	}

	q := application.GetTemplatePertanyaanDefaultByUuidQuery{
		Uuid: uid.String(),
	}

	res, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Nil(t, res)
}
