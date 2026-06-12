package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/common/helper"
	mockaccount "UnpakSiamida/modules/account/application/mock"
	domainaccount "UnpakSiamida/modules/account/domain"
	mockbanksoal "UnpakSiamida/modules/banksoal/application/mock"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateKuesionerCommandHandler(t *testing.T) {
	validBankUUID := uuid.New()
	validDate := "2024-01-01 10:20:30"

	t.Run("Success Simak Dosen", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			CreateFunc: func(ctx context.Context, kuesioner *domainkuesioner.Kuesioner) error {
				assert.Equal(t, "123", kuesioner.NIDN)
				assert.Equal(t, "Dr. John Doe", *kuesioner.NamaDosen)
				assert.Equal(t, "456", kuesioner.IdBankSoal)
				return nil
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				assert.Equal(t, validBankUUID, uid)
				return &domainbanksoal.BankSoal{
					ID: 456,
				}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				assert.Equal(t, "123", *id.NIDN)
				return &domainaccount.AccountDefault{
					ID:   "123",
					Name: helper.StrPtr("Dr. John Doe"),
				}, nil
			},
		}

		handler := &CreateKuesionerCommandHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}

		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Tanggal:      validDate,
			SID:          "123",
			Resource:     "simak",
			CodeCtx:      domainaccount.CtxDosen,
		}

		uid, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, uid)
	})

	t.Run("Success Simak Mahasiswa", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			CreateFunc: func(ctx context.Context, kuesioner *domainkuesioner.Kuesioner) error {
				assert.Equal(t, "npm-789", kuesioner.NPM)
				assert.Equal(t, "Student Alice", *kuesioner.NamaMahasiswa)
				return nil
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{
					ID: 456,
				}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				assert.Equal(t, "npm-789", *id.NIM)
				return &domainaccount.AccountDefault{
					ID:   "npm-789",
					Name: helper.StrPtr("Student Alice"),
				}, nil
			},
		}

		handler := &CreateKuesionerCommandHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}

		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Tanggal:      validDate,
			SID:          "npm-789",
			Resource:     "simak",
			CodeCtx:      domainaccount.CtxMahasiswa,
		}

		uid, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, uid)
	})

	t.Run("Success Simpeg", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			CreateFunc: func(ctx context.Context, kuesioner *domainkuesioner.Kuesioner) error {
				assert.Equal(t, "nip-abc", kuesioner.NIP)
				assert.Equal(t, "Tendik Bob", *kuesioner.NamaTendik)
				return nil
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{
					ID: 456,
				}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				assert.Equal(t, "nip-abc", *id.NIP)
				return &domainaccount.AccountDefault{
					ID:   "nip-abc",
					Name: helper.StrPtr("Tendik Bob"),
				}, nil
			},
		}

		handler := &CreateKuesionerCommandHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}

		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Tanggal:      validDate,
			SID:          "nip-abc",
			Resource:     "simpeg",
		}

		uid, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.NotEmpty(t, uid)
	})

	t.Run("Fail Invalid Bank Soal UUID", func(t *testing.T) {
		handler := &CreateKuesionerCommandHandler{}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: "invalid-uuid",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.InvalidBankSoal(), err)
	})

	t.Run("Fail Bank Soal NotFound", func(t *testing.T) {
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFoundBankSoal(), err)
	})

	t.Run("Fail Bank Soal Other Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return nil, dbErr
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Fail Local Resource RespondentOnly", func(t *testing.T) {
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{ID: 456}, nil
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Resource:     "local",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.RespondentOnly(), err)
	})

	t.Run("Fail Unknown Resource", func(t *testing.T) {
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{ID: 456}, nil
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Resource:     "unknown",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFoundResource(), err)
	})

	t.Run("Fail Account NotFound", func(t *testing.T) {
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{ID: 456}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Resource:     "simpeg",
			SID:          "nip-123",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.NotFoundBankSoal(), err) // returns NotFoundBankSoal on account not found
	})

	t.Run("Fail Account Other Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{ID: 456}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				return nil, dbErr
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Resource:     "simpeg",
			SID:          "nip-123",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Fail NewKuesioner Date Parse Error", func(t *testing.T) {
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{ID: 456}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				return &domainaccount.AccountDefault{ID: "nip-123"}, nil
			},
		}
		handler := &CreateKuesionerCommandHandler{
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Resource:     "simpeg",
			SID:          "nip-123",
			Tanggal:      "invalid-date",
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "period have wrong date format")
	})

	t.Run("Fail Repo Create Error", func(t *testing.T) {
		createErr := errors.New("create err")
		repo := &mockkuesioner.MockKuesionerRepository{
			CreateFunc: func(ctx context.Context, kuesioner *domainkuesioner.Kuesioner) error {
				return createErr
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoal, error) {
				return &domainbanksoal.BankSoal{ID: 456}, nil
			},
		}
		repoAccount := &mockaccount.MockAccountRepository{
			GetFunc: func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
				return &domainaccount.AccountDefault{ID: "nip-123"}, nil
			},
		}
		handler := &CreateKuesionerCommandHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
			RepoAccount:  repoAccount,
		}
		cmd := CreateKuesionerCommand{
			UuidBankSoal: validBankUUID.String(),
			Resource:     "simpeg",
			SID:          "nip-123",
			Tanggal:      validDate,
		}
		_, err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, createErr)
	})
}
