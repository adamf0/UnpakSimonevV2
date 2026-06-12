package mock

import (
	commonDomain "UnpakSiamida/common/domain"
	domainFakultas "UnpakSiamida/modules/fakultas/domain"
	"context"
)

type MockFakultasRepository struct {
	GetAllFunc func(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int) ([]domainFakultas.FakultasDefault, int64, error)
}

var _ domainFakultas.IFakultasRepository = (*MockFakultasRepository)(nil)

func (m *MockFakultasRepository) GetAll(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int) ([]domainFakultas.FakultasDefault, int64, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, search, searchFilters, page, limit)
	}
	return nil, 0, nil
}
