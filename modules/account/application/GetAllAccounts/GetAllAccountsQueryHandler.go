package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainAccount "UnpakSiamida/modules/account/domain"
	"context"
	"time"
)

type GetAllAccountsQueryHandler struct {
	Repo domainAccount.IAccountRepository
}

func (h *GetAllAccountsQueryHandler) Handle(
	ctx context.Context,
	q GetAllAccountsQuery,
) (commondomain.Paged[domainAccount.Account], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Accounts, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
		q.Deleted,
	)
	if err != nil {
		return commondomain.Paged[domainAccount.Account]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainAccount.Account]{
		Data:        Accounts,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
