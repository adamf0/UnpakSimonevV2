package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/common/helper"
	"UnpakSiamida/modules/banksoal/application/mock"
	"UnpakSiamida/modules/banksoal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestScheduleTimeBankSoalCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	t.Run("invalid uuid", func(t *testing.T) {
		repo := &mock.MockRepository{}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: "invalid-uuid",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("bank soal default not found", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return nil, gorm.ErrRecordNotFound
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{UUID: validUUID}, nil
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("bank soal not found", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{UUID: validUUID}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("owner path - invalid date format", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{
					UUID:         validUUID,
					CreatedBy:    helper.StrPtr("local"),
					CreatedByRef: helper.StrPtr("sid-123"),
				}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					UUID: validUUID,
				}, nil
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
			TanggalMulai: "invalid-date",
			TanggalAkhir: "2026-12-31",
			Resource:     "local",
			SID:          "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("owner path - update repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{
					UUID:         validUUID,
					CreatedBy:    helper.StrPtr("local"),
					CreatedByRef: helper.StrPtr("sid-123"),
				}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					UUID: validUUID,
				}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				return errors.New("update error")
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
			TanggalMulai: "2026-06-12",
			TanggalAkhir: "2026-12-31",
			Resource:     "local",
			SID:          "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("owner path - success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{
					UUID:         validUUID,
					CreatedBy:    helper.StrPtr("local"),
					CreatedByRef: helper.StrPtr("sid-123"),
				}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					UUID: validUUID,
				}, nil
			},
			UpdateFunc: func(ctx context.Context, banksoal *domain.BankSoal) error {
				assert.Equal(t, "2026-06-12", banksoal.TanggalMulai.Format("2006-01-02"))
				assert.Equal(t, "2026-12-31", banksoal.TanggalAkhir.Format("2006-01-02"))
				return nil
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
			TanggalMulai: "2026-06-12",
			TanggalAkhir: "2026-12-31",
			Resource:     "local",
			SID:          "sid-123",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, validUUID.String(), res)
	})

	t.Run("ext path - invalid owner (Resource != local)", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{
					UUID:         validUUID,
					CreatedBy:    helper.StrPtr("local"),
					CreatedByRef: helper.StrPtr("sid-owner"),
				}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					UUID: validUUID,
				}, nil
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
			TanggalMulai: "2026-06-12",
			TanggalAkhir: "2026-12-31",
			Resource:     "not-local-ext",
			SID:          "sid-ext",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("ext path - create ext repo error", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{
					UUID:         validUUID,
					CreatedBy:    helper.StrPtr("local"),
					CreatedByRef: helper.StrPtr("sid-owner"),
				}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					UUID: validUUID,
				}, nil
			},
			CreateExtFunc: func(ctx context.Context, banksoalext *domain.BankSoalExt) error {
				return errors.New("insert ext error")
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
			TanggalMulai: "2026-06-12",
			TanggalAkhir: "2026-12-31",
			Resource:     "local",
			SID:          "sid-ext",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("ext path - success", func(t *testing.T) {
		repo := &mock.MockRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
				return &domain.BankSoalDefault{
					UUID:         validUUID,
					CreatedBy:    helper.StrPtr("local"),
					CreatedByRef: helper.StrPtr("sid-owner"),
				}, nil
			},
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
				return &domain.BankSoal{
					UUID: validUUID,
				}, nil
			},
			CreateExtFunc: func(ctx context.Context, banksoalext *domain.BankSoalExt) error {
				assert.Equal(t, "2026-06-12", banksoalext.TanggalMulai.Format("2006-01-02"))
				assert.Equal(t, "2026-12-31", banksoalext.TanggalAkhir.Format("2006-01-02"))
				assert.Equal(t, "local", *banksoalext.CreatedBy)
				assert.Equal(t, "sid-ext", *banksoalext.CreatedByRef)
				return nil
			},
		}
		handler := &ScheduleTimeBankSoalCommandHandler{Repo: repo}
		cmd := ScheduleTimeBankSoalCommand{
			UuidBankSoal: validUUID.String(),
			TanggalMulai: "2026-06-12",
			TanggalAkhir: "2026-12-31",
			Resource:     "local",
			SID:          "sid-ext",
		}
		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
