package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/banksoal/application/mock"
	"UnpakSiamida/modules/banksoal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCreateBankSoalCommandHandler_Handle(t *testing.T) {
	t.Run("invalid owner / resource is not local", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &CreateBankSoalCommandHandler{Repo: repo}
		cmd := CreateBankSoalCommand{
			Judul:     "Test Judul",
			Content:   "Test Content",
			Deskripsi: "Test Deskripsi",
			Semester:  "Gasal",
			Resource:  "not-local",
			SID:       "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("repo create error", func(t *testing.T) {
		repo := &mock.MockRepository{
			CreateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return errors.New("db insert error")
			},
		}
		handler := &CreateBankSoalCommandHandler{Repo: repo}
		cmd := CreateBankSoalCommand{
			Judul:     "Test Judul",
			Content:   "Test Content",
			Deskripsi: "Test Deskripsi",
			Semester:  "Gasal",
			Resource:  "local",
			SID:       "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		repo := &mock.MockRepository{
			CreateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return nil
			},
		}
		handler := &CreateBankSoalCommandHandler{Repo: repo}
		cmd := CreateBankSoalCommand{
			Judul:     "Test Judul",
			Content:   "Test Content",
			Deskripsi: "Test Deskripsi",
			Semester:  "Gasal",
			Resource:  "local",
			SID:       "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
