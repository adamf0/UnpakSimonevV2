package application

import (
	"context"
	"errors"
	"testing"

	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	mockbanksoal "UnpakSiamida/modules/banksoal/application/mock"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestActiveKuesionerSingleQueryHandler(t *testing.T) {
	kuesionerUUID := uuid.New()

	baseQuery := ActiveKuesionerSingleQuery{
		UUID: kuesionerUUID.String(),
		NIDN: helper.StrPtr("nidn-123"),
	}

	t.Run("Fail OnlyStudentLecturerStaff when all identities nil", func(t *testing.T) {
		handler := &ActiveKuesionerSingleQueryHandler{}
		q := ActiveKuesionerSingleQuery{
			UUID: kuesionerUUID.String(),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainbanksoal.OnlyStudentLecturerStaff(), err)
	})

	t.Run("Fail Invalid UUID", func(t *testing.T) {
		handler := &ActiveKuesionerSingleQueryHandler{}
		q := baseQuery
		q.UUID = "invalid-uuid"

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainbanksoal.InvalidUuid(), err)
	})

	t.Run("Fail Repo GetDefaultByUuid Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return nil, dbErr
			},
		}

		handler := &ActiveKuesionerSingleQueryHandler{
			Repo: repo,
		}

		_, err := handler.Handle(context.Background(), baseQuery)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Fail Identity Mismatch", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return &domainkuesioner.KuesionerDefault{
					UUID: kuesionerUUID,
					NIDN: helper.StrPtr("nidn-different"),
				}, nil
			},
		}

		handler := &ActiveKuesionerSingleQueryHandler{
			Repo: repo,
		}

		_, err := handler.Handle(context.Background(), baseQuery)
		assert.Error(t, err)
		assert.Equal(t, "BankSoal.InvalidData", err.(commonDomain.Error).Code)
	})

	t.Run("Fail BankSoal GetDefaultByKuesioner Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return &domainkuesioner.KuesionerDefault{
					UUID: kuesionerUUID,
					NIDN: helper.StrPtr("nidn-123"),
				}, nil
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetDefaultByKuesionerFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoalDefault, error) {
				return nil, dbErr
			},
		}

		handler := &ActiveKuesionerSingleQueryHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
		}

		_, err := handler.Handle(context.Background(), baseQuery)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Fail GetTotalInputByKuesionerIDs Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return &domainkuesioner.KuesionerDefault{
					Id:   100,
					UUID: kuesionerUUID,
					NIDN: helper.StrPtr("nidn-123"),
				}, nil
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetDefaultByKuesionerFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoalDefault, error) {
				return &domainbanksoal.BankSoalDefault{}, nil
			},
		}
		repoJawaban := &mockkuesioner.MockKuesionerJawabanRepository{
			GetTotalInputByKuesionerIDsFunc: func(ctx context.Context, ids []uint) (map[string]uint, error) {
				return nil, dbErr
			},
		}

		handler := &ActiveKuesionerSingleQueryHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
			RepoJawaban:  repoJawaban,
		}

		_, err := handler.Handle(context.Background(), baseQuery)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Success Full Execution", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.KuesionerDefault, error) {
				return &domainkuesioner.KuesionerDefault{
					Id:   100,
					UUID: kuesionerUUID,
					NIDN: helper.StrPtr("nidn-123"),
				}, nil
			},
		}
		repoBankSoal := &mockbanksoal.MockRepository{
			GetDefaultByKuesionerFunc: func(ctx context.Context, uid uuid.UUID) (*domainbanksoal.BankSoalDefault, error) {
				return &domainbanksoal.BankSoalDefault{
					Id: 10,
				}, nil
			},
		}
		repoJawaban := &mockkuesioner.MockKuesionerJawabanRepository{
			GetTotalInputByKuesionerIDsFunc: func(ctx context.Context, ids []uint) (map[string]uint, error) {
				assert.Equal(t, []uint{100}, ids)
				return map[string]uint{
					kuesionerUUID.String(): 4,
				}, nil
			},
		}

		handler := &ActiveKuesionerSingleQueryHandler{
			Repo:         repo,
			RepoBankSoal: repoBankSoal,
			RepoJawaban:  repoJawaban,
		}

		res, err := handler.Handle(context.Background(), baseQuery)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, kuesionerUUID, res.UUIDKuesioner)
		assert.Equal(t, uint(4), res.TotalInput)
	})
}
