package application

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatejawaban/application/mock"
	domaintemplatejawaban "UnpakSiamida/modules/templatejawaban/domain"
	mockTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/mock"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateTemplateJawabanCommandHandler_Handle(t *testing.T) {
	validUUID := uuid.New()
	validPertanyaanUUID := uuid.New()
	validNilaiStr := "10"

	tests := []struct {
		name          string
		cmd           UpdateTemplateJawabanCommand
		mockJawaban   *mock.MockTemplateJawabanRepository
		mockPertanyan *mockTemplatePertanyaan.MockTemplatePertanyaanRepository
		expectedErr   string
	}{
		{
			name: "Success case",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				Jawaban:                "Diubah",
				Nilai:                  &validNilaiStr,
				IsFreeText:             "1",
				SID:                    "sid-update",
				Resource:               "local",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:   1,
						UUID: validUUID,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tj *domaintemplatejawaban.TemplateJawaban) error {
					return nil
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   50,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "",
		},
		{
			name: "Invalid Uuid format",
			cmd: UpdateTemplateJawabanCommand{
				Uuid: "invalid-uuid",
			},
			mockJawaban:   &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{},
			expectedErr:   "TemplateJawaban.InvalidUuid",
		},
		{
			name: "Invalid UuidTemplatePertanyaan format",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: "invalid-uuid",
			},
			mockJawaban:   &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{},
			expectedErr:   "TemplateJawaban.InvalidTemplatePertanyaan",
		},
		{
			name: "Invalid IsFreeText format",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				IsFreeText:             "not-a-number",
			},
			mockJawaban:   &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{},
			expectedErr:   "Must be a positive number",
		},
		{
			name: "Jawaban record not found in DB",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				IsFreeText:             "1",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{},
			expectedErr:   "TemplateJawaban.NotFound",
		},
		{
			name: "Jawaban Repo error",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				IsFreeText:             "1",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return nil, errors.New("read failed")
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{},
			expectedErr:   "read failed",
		},
		{
			name: "Pertanyaan not found in DB",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				IsFreeText:             "1",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:   1,
						UUID: validUUID,
					}, nil
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return nil, errors.New("not found")
				},
			},
			expectedErr: "TemplateJawaban.NotFoundTemplatePertanyaan",
		},
		{
			name: "Domain logic mismatch UUID",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				IsFreeText:             "1",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:   1,
						UUID: uuid.New(), // mismatch uuid
					}, nil
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   50,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "TemplateJawaban.InvalidData",
		},
		{
			name: "Repo Update database failure",
			cmd: UpdateTemplateJawabanCommand{
				Uuid:                   validUUID.String(),
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				IsFreeText:             "1",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatejawaban.TemplateJawaban, error) {
					return &domaintemplatejawaban.TemplateJawaban{
						ID:   1,
						UUID: validUUID,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tj *domaintemplatejawaban.TemplateJawaban) error {
					return errors.New("update query failed")
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   50,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "update query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UpdateTemplateJawabanCommandHandler{
				Repo:           tt.mockJawaban,
				RepoPertanyaan: tt.mockPertanyan,
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
