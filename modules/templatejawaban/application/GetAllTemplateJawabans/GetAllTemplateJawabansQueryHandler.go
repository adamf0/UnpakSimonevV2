package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainTemplateJawaban "UnpakSiamida/modules/templatejawaban/domain"
	"context"
	"time"
)

type GetAllTemplateJawabansQueryHandler struct {
	Repo domainTemplateJawaban.ITemplateJawabanRepository
}

func (h *GetAllTemplateJawabansQueryHandler) Handle(
	ctx context.Context,
	q GetAllTemplateJawabansQuery,
) (commondomain.Paged[domainTemplateJawaban.TemplateJawabanDefault], error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	TemplateJawabans, total, err := h.Repo.GetAll(
		ctx,
		q.Search,
		q.SearchFilters,
		q.Page,
		q.Limit,
		q.Deleted,
	)
	if err != nil {
		return commondomain.Paged[domainTemplateJawaban.TemplateJawabanDefault]{}, err
	}

	currentPage := 1
	totalPages := 1

	if q.Page != nil {
		currentPage = *q.Page
	}
	if q.Limit != nil && *q.Limit > 0 {
		totalPages = int((total + int64(*q.Limit) - 1) / int64(*q.Limit))
	}

	return commondomain.Paged[domainTemplateJawaban.TemplateJawabanDefault]{
		Data:        TemplateJawabans,
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}
