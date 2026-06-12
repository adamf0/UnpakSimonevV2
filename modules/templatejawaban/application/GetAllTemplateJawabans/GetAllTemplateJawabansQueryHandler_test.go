package application

import (
	"context"
	"errors"
	"testing"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/templatejawaban/application/mock"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTemplateJawabansQueryHandler_Handle(t *testing.T) {
	pageVal := 2
	limitVal := 10
	zeroLimit := 0

	tests := []struct {
		name        string
		query       GetAllTemplateJawabansQuery
		mockRepo    *mock.MockTemplateJawabanRepository
		expectedErr string
		verify      func(t *testing.T, res commondomain.Paged[domaintemplatejawaban.TemplateJawabanDefault])
	}{
		{
			name: "Success with pagination",
			query: GetAllTemplateJawabansQuery{
				Search:        "test",
				SearchFilters: []commondomain.SearchFilter{},
				Page:          &pageVal,
				Limit:         &limitVal,
				Deleted:       false,
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domaintemplatejawaban.TemplateJawabanDefault, int64, error) {
					return []domaintemplatejawaban.TemplateJawabanDefault{
						{
							ID:      1,
							UUID:    uuid.New(),
							Jawaban: "Jawaban A",
							Nilai:   5,
						},
						{
							ID:      2,
							UUID:    uuid.New(),
							Jawaban: "Jawaban B",
							Nilai:   4,
						},
					}, 25, nil
				},
			},
			expectedErr: "",
			verify: func(t *testing.T, res commondomain.Paged[domaintemplatejawaban.TemplateJawabanDefault]) {
				assert.Equal(t, int64(25), res.Total)
				assert.Equal(t, 2, res.CurrentPage)
				assert.Equal(t, 3, res.TotalPages) // ceil(25 / 10) = 3
				assert.Len(t, res.Data, 2)
				assert.Equal(t, "Jawaban A", res.Data[0].Jawaban)
			},
		},
		{
			name: "Success with nil page/limit",
			query: GetAllTemplateJawabansQuery{
				Search:        "test",
				SearchFilters: []commondomain.SearchFilter{},
				Page:          nil,
				Limit:         nil,
				Deleted:       false,
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domaintemplatejawaban.TemplateJawabanDefault, int64, error) {
					return []domaintemplatejawaban.TemplateJawabanDefault{
						{
							ID:      1,
							UUID:    uuid.New(),
							Jawaban: "Jawaban A",
						},
					}, 1, nil
				},
			},
			expectedErr: "",
			verify: func(t *testing.T, res commondomain.Paged[domaintemplatejawaban.TemplateJawabanDefault]) {
				assert.Equal(t, int64(1), res.Total)
				assert.Equal(t, 1, res.CurrentPage)
				assert.Equal(t, 1, res.TotalPages)
				assert.Len(t, res.Data, 1)
			},
		},
		{
			name: "Success with zero limit",
			query: GetAllTemplateJawabansQuery{
				Search:        "",
				SearchFilters: nil,
				Page:          &pageVal,
				Limit:         &zeroLimit,
				Deleted:       false,
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domaintemplatejawaban.TemplateJawabanDefault, int64, error) {
					return []domaintemplatejawaban.TemplateJawabanDefault{}, 10, nil
				},
			},
			expectedErr: "",
			verify: func(t *testing.T, res commondomain.Paged[domaintemplatejawaban.TemplateJawabanDefault]) {
				assert.Equal(t, int64(10), res.Total)
				assert.Equal(t, 2, res.CurrentPage)
				assert.Equal(t, 1, res.TotalPages) // totalPages stays 1 if limit is zero or <= 0
			},
		},
		{
			name: "Repo Error",
			query: GetAllTemplateJawabansQuery{
				Search:  "test",
				Deleted: false,
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetAllFunc: func(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, deleted bool) ([]domaintemplatejawaban.TemplateJawabanDefault, int64, error) {
					return nil, 0, errors.New("database error")
				},
			},
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &GetAllTemplateJawabansQueryHandler{
				Repo: tt.mockRepo,
			}
			res, err := handler.Handle(context.Background(), tt.query)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, res)
				}
			}
		})
	}
}
