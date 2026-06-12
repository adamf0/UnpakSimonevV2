package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	domainaccount "UnpakSiamida/modules/account/domain"
	"context"

	"github.com/google/uuid"
)

type MockAccountRepository struct {
	AuthFunc      func(ctx context.Context, username string, password string) (*domainaccount.AccountDefault, error)
	GetFunc       func(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error)
	GetByUuidFunc func(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error)
	GetAllFunc    func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domainaccount.Account, int64, error)
	CreateFunc    func(ctx context.Context, account *domainaccount.Account) error
	UpdateFunc    func(ctx context.Context, account *domainaccount.Account) error
	DeleteFunc    func(ctx context.Context, uid uuid.UUID) error
	SetupUuidFunc func(ctx context.Context) error
}

func (m *MockAccountRepository) Auth(ctx context.Context, username string, password string) (*domainaccount.AccountDefault, error) {
	if m.AuthFunc != nil {
		return m.AuthFunc(ctx, username, password)
	}
	return nil, nil
}

func (m *MockAccountRepository) Get(ctx context.Context, id domainaccount.AccountIdentifier) (*domainaccount.AccountDefault, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAccountRepository) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainaccount.Account, error) {
	if m.GetByUuidFunc != nil {
		return m.GetByUuidFunc(ctx, uid)
	}
	return nil, nil
}

func (m *MockAccountRepository) GetAll(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int, deleted bool) ([]domainaccount.Account, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit, deleted)
	}
	return nil, 0, nil
}

func (m *MockAccountRepository) Create(ctx context.Context, account *domainaccount.Account) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, account)
	}
	return nil
}

func (m *MockAccountRepository) Update(ctx context.Context, account *domainaccount.Account) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, account)
	}
	return nil
}

func (m *MockAccountRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uid)
	}
	return nil
}

func (m *MockAccountRepository) SetupUuid(ctx context.Context) error {
	if m.SetupUuidFunc != nil {
		return m.SetupUuidFunc(ctx)
	}
	return nil
}
