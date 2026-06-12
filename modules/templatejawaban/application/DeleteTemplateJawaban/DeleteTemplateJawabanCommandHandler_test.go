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

func TestDeleteTemplateJawabanCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name        string
		cmd         DeleteTemplateJawabanCommand
		mockRepo    *mock.MockTemplateJawabanRepository
		expectedErr string
	}{
		{
			name: "Success Hard Delete",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "hard_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
					return nil
				},
			},
			expectedErr: "",
		},
		{
			name: "Success Soft Delete",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "soft_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:   1,
						UUID: validUUID,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tj *domaintemplatejawaban.TemplateJawaban) error {
					assert.NotNil(t, tj.DeletedAt)
					return nil
				},
			},
			expectedErr: "",
		},
		{
			name: "Invalid UUID",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: "invalid-uuid",
				Mode: "hard_delete",
			},
			mockRepo:    &mock.MockTemplateJawabanRepository{},
			expectedErr: "TemplateJawaban.InvalidUuid",
		},
		{
			name: "Record Not Found",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "soft_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			expectedErr: "TemplateJawaban.NotFound",
		},
		{
			name: "GetByUuid Other Error",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "soft_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, errors.New("db error")
				},
			},
			expectedErr: "db error",
		},
		{
			name: "Hard Delete Repo Failure",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "hard_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				DeleteFunc: func(ctx context.Context, uid uuid.UUID) error {
					return errors.New("delete error")
				},
			},
			expectedErr: "delete error",
		},
		{
			name: "Soft Delete Domain Failure (Nil Prev)",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "soft_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, nil // causes DeleteTemplateJawaban to return EmptyData failure
				},
			},
			expectedErr: "TemplateJawaban.EmptyData",
		},
		{
			name: "Soft Delete Repo Update Failure",
			cmd: DeleteTemplateJawabanCommand{
				Uuid: validUUID.String(),
				Mode: "soft_delete",
			},
			mockRepo: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:   1,
						UUID: validUUID,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tj *domaintemplatejawaban.TemplateJawaban) error {
					return errors.New("update error")
				},
			},
			expectedErr: "update error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &DeleteTemplateJawabanCommandHandler{
				Repo: tt.mockRepo,
			}
			res, err := handler.Handle(context.Background(), tt.cmd)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.cmd.Uuid, res)
			}
		})
	}
}
