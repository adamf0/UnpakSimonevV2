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

func TestGetTemplateJawabanByUuidQueryHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name        string
		query       GetTemplateJawabanByUuidQuery
		mockRepo    *mock.MockTemplateJawabanRepository
		expectedErr string
		expectNil   bool
	}{
		{
			name: "Success case",
			query: GetTemplateJawabanByUuidQuery{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:      101,
						UUID:    validUUID,
						Jawaban: "Baik",
					}, nil
				},
			},
			expectedErr: "",
			expectNil:   false,
		},
		{
			name: "Invalid UUID format",
			query: GetTemplateJawabanByUuidQuery{
				Uuid: "invalid-uuid",
			},
			mockRepo:    &mock.MockTemplateJawabanRepository{},
			expectedErr: "TemplateJawaban.NotFound",
			expectNil:   true,
		},
		{
			name: "Not Found (gorm.ErrRecordNotFound)",
			query: GetTemplateJawabanByUuidQuery{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			expectedErr: "TemplateJawaban.NotFound",
			expectNil:   true,
		},
		{
			name: "Repository Other Error",
			query: GetTemplateJawabanByUuidQuery{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, errors.New("connection timed out")
				},
			},
			expectedErr: "connection timed out",
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &GetTemplateJawabanByUuidQueryHandler{
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
				assert.Equal(t, "Baik", res.Jawaban)
			}
		})
	}
}
