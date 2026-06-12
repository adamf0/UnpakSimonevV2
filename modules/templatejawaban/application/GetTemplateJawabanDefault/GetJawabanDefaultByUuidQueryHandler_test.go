package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatejawaban/application/mock"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetTemplateJawabanDefaultByUuidQueryHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name        string
		query       GetTemplateJawabanDefaultByUuidQuery
		mockRepo    *mock.MockTemplateJawabanRepository
		expectedErr string
		expectNil   bool
	}{
		{
			name: "Success case",
			query: GetTemplateJawabanDefaultByUuidQuery{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawabanDefault, error) {
					return &domaintemplatejawaban.TemplateJawabanDefault{
						ID:      202,
						UUID:    validUUID,
						Jawaban: "Default Jawaban",
					}, nil
				},
			},
			expectedErr: "",
			expectNil:   false,
		},
		{
			name: "Invalid UUID format",
			query: GetTemplateJawabanDefaultByUuidQuery{
				Uuid: "invalid-uuid",
			},
			mockRepo:    &mock.MockTemplateJawabanRepository{},
			expectedErr: "TemplateJawaban.NotFound",
			expectNil:   true,
		},
		{
			name: "Not Found (gorm.ErrRecordNotFound)",
			query: GetTemplateJawabanDefaultByUuidQuery{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawabanDefault, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			expectedErr: "TemplateJawaban.NotFound",
			expectNil:   true,
		},
		{
			name: "Repository Other Error",
			query: GetTemplateJawabanDefaultByUuidQuery{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetDefaultByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawabanDefault, error) {
					return nil, errors.New("query error")
				},
			},
			expectedErr: "query error",
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &GetTemplateJawabanDefaultByUuidQueryHandler{
				Repo: tt.mockRepo,
			}
			res, err := handler.Handle(context.Background(), tt.query)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				if tt.expectNil {
					assert.Nil(t, res)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, validUUID, res.UUID)
				assert.Equal(t, "Default Jawaban", res.Jawaban)
			}
		})
	}
}
