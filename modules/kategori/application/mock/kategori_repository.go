package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type MockKategoriRepository struct {
	CountCopyFunc         func(ctx context.Context, judul string) (int, error)
	GetByUuidFunc         func(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error)
	GetDefaultByUuidFunc  func(ctx context.Context, uid uuid.UUID) (*domainkategori.KategoriDefault, error)
	GetChildrenFunc       func(ctx context.Context, parentID int) ([]domainkategori.Kategori, error)
	RebuildFullTextFunc   func(ctx context.Context) error
	GetAllFunc            func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domainkategori.KategoriDefault, int64, error)
	CreateFunc            func(ctx context.Context, kategori *domainkategori.Kategori) error
	UpdateFunc            func(ctx context.Context, kategori *domainkategori.Kategori) error
	DeleteFunc            func(ctx context.Context, uid uuid.UUID) error
	SetupUuidFunc         func(ctx context.Context) error
	WithTxFunc            func(tx any) domainkategori.IKategoriRepository
	BeginTxFunc           func(ctx context.Context) (*gorm.DB, error)
	UpdateParentBatchFunc func(ctx context.Context, rows []domainkategori.UpdateRow) error
}

// Ensure it implements the interface
var _ domainkategori.IKategoriRepository = (*MockKategoriRepository)(nil)

func (m *MockKategoriRepository) CountCopy(ctx context.Context, judul string) (int, error) {
	if m.CountCopyFunc != nil {
		return m.CountCopyFunc(ctx, judul)
	}
	return 0, nil
}

func (m *MockKategoriRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainkategori.Kategori, error) {
	if m.GetByUuidFunc != nil {
		return m.GetByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockKategoriRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainkategori.KategoriDefault, error) {
	if m.GetDefaultByUuidFunc != nil {
		return m.GetDefaultByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockKategoriRepository) GetChildren(ctx context.Context, parentID int) ([]domainkategori.Kategori, error) {
	if m.GetChildrenFunc != nil {
		return m.GetChildrenFunc(ctx, parentID)
	}
	return nil, nil
}

func (m *MockKategoriRepository) RebuildFullText(ctx context.Context) error {
	if m.RebuildFullTextFunc != nil {
		return m.RebuildFullTextFunc(ctx)
	}
	return nil
}

func (m *MockKategoriRepository) GetAll(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domainkategori.KategoriDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit, deleted)
	}
	return nil, 0, nil
}

func (m *MockKategoriRepository) Create(ctx context.Context, kategori *domainkategori.Kategori) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, kategori)
	}
	return nil
}

func (m *MockKategoriRepository) Update(ctx context.Context, kategori *domainkategori.Kategori) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, kategori)
	}
	return nil
}

func (m *MockKategoriRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uid)
	}
	return nil
}

func (m *MockKategoriRepository) SetupUuid(ctx context.Context) error {
	if m.SetupUuidFunc != nil {
		return m.SetupUuidFunc(ctx)
	}
	return nil
}

func (m *MockKategoriRepository) WithTx(tx any) domainkategori.IKategoriRepository {
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

func (m *MockKategoriRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{DryRun: true})
	tx := db.Session(&gorm.Session{})
	tx.Statement.ConnPool = &dummyTxConn{ConnPool: db.Statement.ConnPool}
	return tx, nil
}

func (m *MockKategoriRepository) UpdateParentBatch(ctx context.Context, rows []domainkategori.UpdateRow) error {
	if m.UpdateParentBatchFunc != nil {
		return m.UpdateParentBatchFunc(ctx, rows)
	}
	return nil
}
