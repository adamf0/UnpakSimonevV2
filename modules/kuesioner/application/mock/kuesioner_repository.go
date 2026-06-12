package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/kuesioner/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type MockKuesionerRepository struct {
	GetAllKuesionerResultFunc        func(ctx context.Context, JudulBankSoal *string, Is4Year bool, PartitionKey string) ([]domain.KuesionerResult, error)
	GetByUuidFunc                    func(ctx context.Context, uid uuid.UUID) (*domain.Kuesioner, error)
	GetDefaultByUuidFunc             func(ctx context.Context, uid uuid.UUID) (*domain.KuesionerDefault, error)
	GetAllFormFromActiveBankSoalFunc func(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]domain.KuesionerDefault, error)
	GetAllFunc                       func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domain.KuesionerDefault, int64, error)
	CreateFunc                       func(ctx context.Context, kuesioner *domain.Kuesioner) error
	DeleteFunc                       func(ctx context.Context, uid uuid.UUID) error
	SetupUuidFunc                    func(ctx context.Context) error
}

var _ domain.IKuesionerRepository = (*MockKuesionerRepository)(nil)

func (m *MockKuesionerRepository) GetAllKuesionerResult(ctx context.Context, JudulBankSoal *string, Is4Year bool, PartitionKey string) ([]domain.KuesionerResult, error) {
	if m.GetAllKuesionerResultFunc != nil {
		return m.GetAllKuesionerResultFunc(ctx, JudulBankSoal, Is4Year, PartitionKey)
	}
	return nil, nil
}

func (m *MockKuesionerRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.Kuesioner, error) {
	if m.GetByUuidFunc != nil {
		return m.GetByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockKuesionerRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domain.KuesionerDefault, error) {
	if m.GetDefaultByUuidFunc != nil {
		return m.GetDefaultByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockKuesionerRepository) GetAllFormFromActiveBankSoal(ctx context.Context, nidn string, nip string, npm string, banksoal []uint) ([]domain.KuesionerDefault, error) {
	if m.GetAllFormFromActiveBankSoalFunc != nil {
		return m.GetAllFormFromActiveBankSoalFunc(ctx, nidn, nip, npm, banksoal)
	}
	return nil, nil
}

func (m *MockKuesionerRepository) GetAll(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domain.KuesionerDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit, deleted)
	}
	return nil, 0, nil
}

func (m *MockKuesionerRepository) Create(ctx context.Context, kuesioner *domain.Kuesioner) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, kuesioner)
	}
	return nil
}

func (m *MockKuesionerRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uid)
	}
	return nil
}

func (m *MockKuesionerRepository) SetupUuid(ctx context.Context) error {
	if m.SetupUuidFunc != nil {
		return m.SetupUuidFunc(ctx)
	}
	return nil
}

type MockKuesionerJawabanRepository struct {
	GetByPertanyaanAndUserFunc      func(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domain.KuesionerJawaban, error)
	GetTotalInputByKuesionerIDsFunc func(ctx context.Context, ids []uint) (map[string]uint, error)
	GetAllByKuesionerFunc           func(ctx context.Context, uuidkuesioner string) ([]domain.KuesionerJawabanDefault, error)
	CreateFunc                      func(ctx context.Context, data *domain.KuesionerJawaban) error
	DeleteFunc                      func(ctx context.Context, id uint) error
	WithTxFunc                      func(tx any) domain.IKuesionerJawabanRepository
	BeginTxFunc                     func(ctx context.Context) (*gorm.DB, error)
}

var _ domain.IKuesionerJawabanRepository = (*MockKuesionerJawabanRepository)(nil)

func (m *MockKuesionerJawabanRepository) GetByPertanyaanAndUser(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domain.KuesionerJawaban, error) {
	if m.GetByPertanyaanAndUserFunc != nil {
		return m.GetByPertanyaanAndUserFunc(ctx, pertanyaanID, sid, resource)
	}
	return nil, nil
}

func (m *MockKuesionerJawabanRepository) GetTotalInputByKuesionerIDs(ctx context.Context, ids []uint) (map[string]uint, error) {
	if m.GetTotalInputByKuesionerIDsFunc != nil {
		return m.GetTotalInputByKuesionerIDsFunc(ctx, ids)
	}
	return nil, nil
}

func (m *MockKuesionerJawabanRepository) GetAllByKuesioner(ctx context.Context, uuidkuesioner string) ([]domain.KuesionerJawabanDefault, error) {
	if m.GetAllByKuesionerFunc != nil {
		return m.GetAllByKuesionerFunc(ctx, uuidkuesioner)
	}
	return nil, nil
}

func (m *MockKuesionerJawabanRepository) Create(ctx context.Context, data *domain.KuesionerJawaban) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, data)
	}
	return nil
}

func (m *MockKuesionerJawabanRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockKuesionerJawabanRepository) WithTx(tx any) domain.IKuesionerJawabanRepository {
	if m.WithTxFunc != nil {
		return m.WithTxFunc(tx)
	}
	return m
}

type dummyTxConn struct {
	gorm.ConnPool
}

func (dummyTxConn) Commit() error {
	return nil
}

func (dummyTxConn) Rollback() error {
	return nil
}

func (m *MockKuesionerJawabanRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{DryRun: true})
	tx := db.Session(&gorm.Session{})
	tx.Statement.ConnPool = &dummyTxConn{ConnPool: db.Statement.ConnPool}
	return tx, nil
}
