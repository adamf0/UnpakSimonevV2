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

func TestActiveKuesionerQueryHandler(t *testing.T) {
	t.Run("Fail OnlyStudentLecturerStaff when all identity fields nil", func(t *testing.T) {
		handler := &ActiveKuesionerQueryHandler{}
		q := ActiveKuesionerQuery{
			NIDN: nil,
			NIP:  nil,
			NPM:  nil,
		}

		_, err := handler.Handle(context.Background(), q)
		assert.Error(t, err)
		assert.Equal(t, domainbanksoal.OnlyStudentLecturerStaff(), err)
	})

	t.Run("Fail GetAll BankSoal Returns Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repoBankSoal := &mockbanksoal.MockRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, targetFakultas, targetProdi, targetUnit, targetStatus string, page, limit *int, deleted, active bool) ([]domainbanksoal.BankSoalDefault, int64, error) {
				return nil, 0, dbErr
			},
		}

		handler := &ActiveKuesionerQueryHandler{
			RepoBankSoal: repoBankSoal,
		}

		q := ActiveKuesionerQuery{
			NIDN: helper.StrPtr("nidn-123"),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Success Empty Active BankSoal", func(t *testing.T) {
		repoBankSoal := &mockbanksoal.MockRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, targetFakultas, targetProdi, targetUnit, targetStatus string, page, limit *int, deleted, active bool) ([]domainbanksoal.BankSoalDefault, int64, error) {
				return []domainbanksoal.BankSoalDefault{}, 0, nil
			},
		}

		handler := &ActiveKuesionerQueryHandler{
			RepoBankSoal: repoBankSoal,
		}

		q := ActiveKuesionerQuery{
			NIDN: helper.StrPtr("nidn-123"),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("Fail GetAllFormFromActiveBankSoal Returns Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repoBankSoal := &mockbanksoal.MockRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, targetFakultas, targetProdi, targetUnit, targetStatus string, page, limit *int, deleted, active bool) ([]domainbanksoal.BankSoalDefault, int64, error) {
				return []domainbanksoal.BankSoalDefault{
					{
						Id: 10,
					},
				}, 1, nil
			},
		}
		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllFormFromActiveBankSoalFunc: func(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]domainkuesioner.KuesionerDefault, error) {
				return nil, dbErr
			},
		}

		handler := &ActiveKuesionerQueryHandler{
			RepoBankSoal: repoBankSoal,
			Repo:         repo,
		}

		q := ActiveKuesionerQuery{
			NIDN: helper.StrPtr("nidn-123"),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Fail GetTotalInputByKuesionerIDs Returns Error", func(t *testing.T) {
		dbErr := errors.New("db error")
		repoBankSoal := &mockbanksoal.MockRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, targetFakultas, targetProdi, targetUnit, targetStatus string, page, limit *int, deleted, active bool) ([]domainbanksoal.BankSoalDefault, int64, error) {
				return []domainbanksoal.BankSoalDefault{
					{
						Id: 10,
					},
				}, 1, nil
			},
		}
		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllFormFromActiveBankSoalFunc: func(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]domainkuesioner.KuesionerDefault, error) {
				return []domainkuesioner.KuesionerDefault{
					{
						Id:         100,
						UUID:       uuid.New(),
						IdBankSoal: "10",
					},
				}, nil
			},
		}
		repoJawaban := &mockkuesioner.MockKuesionerJawabanRepository{
			GetTotalInputByKuesionerIDsFunc: func(ctx context.Context, ids []uint) (map[string]uint, error) {
				return nil, dbErr
			},
		}

		handler := &ActiveKuesionerQueryHandler{
			RepoBankSoal: repoBankSoal,
			Repo:         repo,
			RepoJawaban:  repoJawaban,
		}

		q := ActiveKuesionerQuery{
			NIDN: helper.StrPtr("nidn-123"),
		}

		_, err := handler.Handle(context.Background(), q)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("Success Full Execution", func(t *testing.T) {
		kuesionerUUID := uuid.New()
		repoBankSoal := &mockbanksoal.MockRepository{
			GetAllFunc: func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, targetFakultas, targetProdi, targetUnit, targetStatus string, page, limit *int, deleted, active bool) ([]domainbanksoal.BankSoalDefault, int64, error) {
				assert.Equal(t, "active", targetStatus)
				return []domainbanksoal.BankSoalDefault{
					{
						Id: 10,
					},
				}, 1, nil
			},
		}
		repo := &mockkuesioner.MockKuesionerRepository{
			GetAllFormFromActiveBankSoalFunc: func(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]domainkuesioner.KuesionerDefault, error) {
				assert.Equal(t, "nidn-123", nidn)
				assert.Equal(t, []uint{10}, banksoal)
				return []domainkuesioner.KuesionerDefault{
					{
						Id:         100,
						UUID:       kuesionerUUID,
						IdBankSoal: "10",
					},
				}, nil
			},
		}
		repoJawaban := &mockkuesioner.MockKuesionerJawabanRepository{
			GetTotalInputByKuesionerIDsFunc: func(ctx context.Context, ids []uint) (map[string]uint, error) {
				assert.Equal(t, []uint{100}, ids)
				return map[string]uint{
					kuesionerUUID.String(): 5,
				}, nil
			},
		}

		handler := &ActiveKuesionerQueryHandler{
			RepoBankSoal: repoBankSoal,
			Repo:         repo,
			RepoJawaban:  repoJawaban,
		}

		q := ActiveKuesionerQuery{
			NIDN: helper.StrPtr("nidn-123"),
		}

		res, err := handler.Handle(context.Background(), q)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, kuesionerUUID, res[0].UUIDKuesioner)
		assert.Equal(t, uint(5), res[0].TotalInput)
	})
}
