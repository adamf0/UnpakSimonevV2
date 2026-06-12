package mock

import (
	"context"
	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type MockTemplatePertanyaanRepository struct {
	CountCopyFunc                       func(ctx context.Context, judul string) (int, error)
	GetByUuidFunc                       func(ctx context.Context, uid uuid.UUID) (*domain.TemplatePertanyaan, error)
	GetDefaultByUuidFunc                  func(ctx context.Context, uid uuid.UUID) (*domain.TemplatePertanyaanDefault, error)
	GetDefaultWithAnswareByUuidFunc       func(ctx context.Context, uid uuid.UUID) (*domain.TemplatePertanyaanWithAnswareDefault, error)
	GetDefaultWithAnswareByBankSoalFunc   func(ctx context.Context, id_banksoal uint) ([]domain.TemplatePertanyaanWithAnswareDefault, error)
	GetAllFunc                            func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domain.TemplatePertanyaanDefault, int64, error)
	CopyByBankSoalFunc                    func(ctx context.Context, tx *gorm.DB, sourceBankSoalID uint, targetBankSoalID uint, resource string, sid string) (map[uint]uint, error)
	CreateFunc                            func(ctx context.Context, aktivitasproker *domain.TemplatePertanyaan) error
	UpdateFunc                            func(ctx context.Context, aktivitasproker *domain.TemplatePertanyaan) error
	DeleteFunc                            func(ctx context.Context, uid uuid.UUID) error
	SetupUuidFunc                         func(ctx context.Context) error
	WithTxFunc                            func(tx any) domain.ITemplatePertanyaanRepository
	BeginTxFunc                           func(ctx context.Context) (*gorm.DB, error)
}

func (m *MockTemplatePertanyaanRepository) CountCopy(ctx context.Context, judul string) (int, error) {
	if m.CountCopyFunc != nil {
		return m.CountCopyFunc(ctx, judul)
	}
	return 0, nil
}

func (m *MockTemplatePertanyaanRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.TemplatePertanyaan, error) {
	if m.GetByUuidFunc != nil {
		return m.GetByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockTemplatePertanyaanRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domain.TemplatePertanyaanDefault, error) {
	if m.GetDefaultByUuidFunc != nil {
		return m.GetDefaultByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockTemplatePertanyaanRepository) GetDefaultWithAnswareByUuid(ctx context.Context, uid uuid.UUID) (*domain.TemplatePertanyaanWithAnswareDefault, error) {
	if m.GetDefaultWithAnswareByUuidFunc != nil {
		return m.GetDefaultWithAnswareByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockTemplatePertanyaanRepository) GetDefaultWithAnswareByBankSoal(ctx context.Context, id_banksoal uint) ([]domain.TemplatePertanyaanWithAnswareDefault, error) {
	if m.GetDefaultWithAnswareByBankSoalFunc != nil {
		return m.GetDefaultWithAnswareByBankSoalFunc(ctx, id_banksoal)
	}
	return nil, nil
}

func (m *MockTemplatePertanyaanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commonDomain.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domain.TemplatePertanyaanDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit, deleted)
	}
	return nil, 0, nil
}

func (m *MockTemplatePertanyaanRepository) CopyByBankSoal(
	ctx context.Context,
	tx *gorm.DB,
	sourceBankSoalID uint,
	targetBankSoalID uint,
	resource string,
	sid string,
) (map[uint]uint, error) {
	if m.CopyByBankSoalFunc != nil {
		return m.CopyByBankSoalFunc(ctx, tx, sourceBankSoalID, targetBankSoalID, resource, sid)
	}
	return nil, nil
}

func (m *MockTemplatePertanyaanRepository) Create(ctx context.Context, aktivitasproker *domain.TemplatePertanyaan) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, aktivitasproker)
	}
	return nil
}

func (m *MockTemplatePertanyaanRepository) Update(ctx context.Context, aktivitasproker *domain.TemplatePertanyaan) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, aktivitasproker)
	}
	return nil
}

func (m *MockTemplatePertanyaanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uid)
	}
	return nil
}

func (m *MockTemplatePertanyaanRepository) SetupUuid(ctx context.Context) error {
	if m.SetupUuidFunc != nil {
		return m.SetupUuidFunc(ctx)
	}
	return nil
}

func (m *MockTemplatePertanyaanRepository) WithTx(tx any) domain.ITemplatePertanyaanRepository {
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

func (m *MockTemplatePertanyaanRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{DryRun: true})
	tx := db.Session(&gorm.Session{})
	tx.Statement.ConnPool = &dummyTxConn{ConnPool: db.Statement.ConnPool}
	return tx, nil
}
