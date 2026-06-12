package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"UnpakSiamida/modules/templatejawaban/application/mock"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRestoreTemplateJawabanCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()
	deletedTime := time.Now()

	tests := []struct {
		name        string
		cmd         RestoreTemplateJawabanCommand
		mockRepo    *mock.MockTemplateJawabanRepository
		expectedErr string
	}{
		{
			name: "Success case",
			cmd: RestoreTemplateJawabanCommand{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:        123,
						UUID:      validUUID,
						DeletedAt: &deletedTime,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tj *domaintemplatejawaban.TemplateJawaban) error {
					assert.Nil(t, tj.DeletedAt)
					return nil
				},
			},
			expectedErr: "",
		},
		{
			name: "Invalid UUID format",
			cmd: RestoreTemplateJawabanCommand{
				Uuid: "invalid-uuid",
			},
			mockRepo:    &mock.MockTemplateJawabanRepository{},
			expectedErr: "TemplateJawaban.InvalidUuid",
		},
		{
			name: "Not Found (gorm.ErrRecordNotFound)",
			cmd: RestoreTemplateJawabanCommand{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			expectedErr: "TemplateJawaban.NotFound",
		},
		{
			name: "Repository GetByUuid Other Error",
			cmd: RestoreTemplateJawabanCommand{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, errors.New("read error")
				},
			},
			expectedErr: "read error",
		},
		{
			name: "Domain Restore Failure (Nil Prev)",
			cmd: RestoreTemplateJawabanCommand{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, nil
				},
			},
			expectedErr: "TemplateJawaban.EmptyData",
		},
		{
			name: "Repository Update Error",
			cmd: RestoreTemplateJawabanCommand{
				Uuid: validUUID.String(),
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:        123,
						UUID:      validUUID,
						DeletedAt: &deletedTime,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tj *domaintemplatejawaban.TemplateJawaban) error {
					return errors.New("write error")
				},
			},
			expectedErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &RestoreTemplateJawabanCommandHandler{
				Repo: tt.mockRepo,
			}
			res, err := handler.Handle(context.Background(), tt.cmd)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, validUUID.String(), res)
			}
		})
	}
}
