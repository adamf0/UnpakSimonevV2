package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	domainProdi "UnpakSiamida/modules/prodi/domain"
	"context"
)

type MockProdiRepository struct {
	GetAllFunc func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int) ([]domainProdi.ProdiDefault, int64, error)
}

var _ domainProdi.IProdiRepository = (*MockProdiRepository)(nil)

func (m *MockProdiRepository) GetAll(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int) ([]domainProdi.ProdiDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit)
	}
	return nil, 0, nil
}
