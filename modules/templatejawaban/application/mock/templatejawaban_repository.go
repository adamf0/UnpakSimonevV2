package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/templatejawaban/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type MockTemplateJawabanRepository struct {
	GetByUuidFunc                func(ctx context.Context, uid uuid.UUID) (*domain.TemplateJawaban, error)
	GetByUUIDsFunc               func(ctx context.Context, uuids []string) ([]domain.TemplateJawaban, error)
	GetFreeTextByPertanyaanFunc  func(ctx context.Context, pertanyaanID uint) (*domain.TemplateJawaban, error)
	GetDefaultByUuidFunc         func(ctx context.Context, uid uuid.UUID) (*domain.TemplateJawabanDefault, error)
	GetAllFunc                   func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domain.TemplateJawabanDefault, int64, error)
	CreateFunc                   func(ctx context.Context, templatejawaban *domain.TemplateJawaban) error
	UpdateFunc                   func(ctx context.Context, templatejawaban *domain.TemplateJawaban) error
	DeleteFunc                   func(ctx context.Context, uid uuid.UUID) error
	CopyByTemplatePertanyaanFunc func(ctx context.Context, tx *gorm.DB, sourceTemplatePertanyaanID uint, targetTemplatePertanyaanID uint) error
	SetupUuidFunc                func(ctx context.Context) error
	WithTxFunc                   func(tx any) domain.ITemplateJawabanRepository
	BeginTxFunc                  func(ctx context.Context) (*gorm.DB, error)
}

// Ensure it implements the interface
var _ domain.ITemplateJawabanRepository = (*MockTemplateJawabanRepository)(nil)

func (m *MockTemplateJawabanRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.TemplateJawaban, error) {
	if m.GetByUuidFunc != nil {
		return m.GetByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockTemplateJawabanRepository) GetByUUIDs(ctx context.Context, uuids []string) ([]domain.TemplateJawaban, error) {
	if m.GetByUUIDsFunc != nil {
		return m.GetByUUIDsFunc(ctx, uuids)
	}
	return nil, nil
}

func (m *MockTemplateJawabanRepository) GetFreeTextByPertanyaan(ctx context.Context, pertanyaanID uint) (*domain.TemplateJawaban, error) {
	if m.GetFreeTextByPertanyaanFunc != nil {
		return m.GetFreeTextByPertanyaanFunc(ctx, pertanyaanID)
	}
	return nil, nil
}

func (m *MockTemplateJawabanRepository) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domain.TemplateJawabanDefault, error) {
	if m.GetDefaultByUuidFunc != nil {
		return m.GetDefaultByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockTemplateJawabanRepository) GetAll(
	ctx context.Context,
	search string,
	searchFilters []commonDomain.SearchFilter,
	page, limit *int,
	deleted bool,
) ([]domain.TemplateJawabanDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit, deleted)
	}
	return nil, 0, nil
}

func (m *MockTemplateJawabanRepository) Create(ctx context.Context, templatejawaban *domain.TemplateJawaban) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, templatejawaban)
	}
	return nil
}

func (m *MockTemplateJawabanRepository) Update(ctx context.Context, templatejawaban *domain.TemplateJawaban) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, templatejawaban)
	}
	return nil
}

func (m *MockTemplateJawabanRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uid)
	}
	return nil
}

func (m *MockTemplateJawabanRepository) CopyByTemplatePertanyaan(
	ctx context.Context,
	tx *gorm.DB,
	sourceTemplatePertanyaanID uint,
	targetTemplatePertanyaanID uint,
) error {
	if m.CopyByTemplatePertanyaanFunc != nil {
		return m.CopyByTemplatePertanyaanFunc(ctx, tx, sourceTemplatePertanyaanID, targetTemplatePertanyaanID)
	}
	return nil
}

func (m *MockTemplateJawabanRepository) SetupUuid(ctx context.Context) error {
	if m.SetupUuidFunc != nil {
		return m.SetupUuidFunc(ctx)
	}
	return nil
}

func (m *MockTemplateJawabanRepository) WithTx(tx any) domain.ITemplateJawabanRepository {
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

func (m *MockTemplateJawabanRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{DryRun: true})
	tx := db.Session(&gorm.Session{})
	tx.Statement.ConnPool = &dummyTxConn{ConnPool: db.Statement.ConnPool}
	return tx, nil
}
