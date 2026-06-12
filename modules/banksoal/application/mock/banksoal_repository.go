package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/banksoal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type MockRepository struct {
	CountCopyFunc             func(ctx context.Context, judul string) (int, error)
	GetByUuidFunc             func(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error)
	GetDefaultByUuidFunc      func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error)
	GetDefaultByKuesionerFunc func(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error)
	GetAllFunc                func(
		ctx context.Context,
		search string,
		searchFilters []commonDomain.SearchFilter,
		TargetFakultas string,
		TargetProdi string,
		TargetUnit string,
		TargetStatus string,
		page, limit *int,
		deleted bool,
		active bool,
	) ([]domain.BankSoalDefault, int64, error)
	CreateFunc    func(ctx context.Context, banksoal *domain.BankSoal) error
	UpdateFunc    func(ctx context.Context, banksoal *domain.BankSoal) error
	DeleteFunc    func(ctx context.Context, uid uuid.UUID) error
	SetupUuidFunc func(ctx context.Context) error

	CreateExtFunc func(ctx context.Context, banksoalext *domain.BankSoalExt) error
	DeleteExtFunc func(ctx context.Context, uid uuid.UUID, idbanksoal uint) error

	WithTxFunc  func(tx any) domain.IBankSoalRepository
	BeginTxFunc func(ctx context.Context) (*gorm.DB, error)
}

var _ domain.IBankSoalRepository = (*MockRepository)(nil)

func (m *MockRepository) CountCopy(ctx context.Context, judul string) (int, error) {
	if m.CountCopyFunc != nil {
		return m.CountCopyFunc(ctx, judul)
	}
	return 0, nil
}

func (m *MockRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.BankSoal, error) {
	if m.GetByUuidFunc != nil {
		return m.GetByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
	if m.GetDefaultByUuidFunc != nil {
		return m.GetDefaultByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockRepository) GetDefaultByKuesioner(ctx context.Context, uid uuid.UUID) (*domain.BankSoalDefault, error) {
	if m.GetDefaultByKuesionerFunc != nil {
		return m.GetDefaultByKuesionerFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commonDomain.SearchFilter,
	TargetFakultas string,
	TargetProdi string,
	TargetUnit string,
	TargetStatus string,
	page, limit *int,
	deleted bool,
	active bool,
) ([]domain.BankSoalDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, TargetFakultas, TargetProdi, TargetUnit, TargetStatus, page, limit, deleted, active)
	}
	return nil, 0, nil
}

func (m *MockRepository) Create(ctx context.Context, banksoal *domain.BankSoal) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, banksoal)
	}
	return nil
}

func (m *MockRepository) Update(ctx context.Context, banksoal *domain.BankSoal) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, banksoal)
	}
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uid)
	}
	return nil
}

func (m *MockRepository) SetupUuid(ctx context.Context) error {
	if m.SetupUuidFunc != nil {
		return m.SetupUuidFunc(ctx)
	}
	return nil
}

func (m *MockRepository) CreateExt(ctx context.Context, banksoalext *domain.BankSoalExt) error {
	if m.CreateExtFunc != nil {
		return m.CreateExtFunc(ctx, banksoalext)
	}
	return nil
}

func (m *MockRepository) DeleteExt(ctx context.Context, uid uuid.UUID, idbanksoal uint) error {
	if m.DeleteExtFunc != nil {
		return m.DeleteExtFunc(ctx, uid, idbanksoal)
	}
	return nil
}

func (m *MockRepository) WithTx(tx any) domain.IBankSoalRepository {
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

func (m *MockRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{DryRun: true})
	tx := db.Session(&gorm.Session{})
	tx.Statement.ConnPool = &dummyTxConn{ConnPool: db.Statement.ConnPool}
	return tx, nil
}
