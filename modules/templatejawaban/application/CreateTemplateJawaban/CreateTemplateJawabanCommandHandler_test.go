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
)

func TestCreateTemplateJawabanCommandHandler_Handle(t *testing.T) {
	validPertanyaanUUID := uuid.New()
	validNilaiStr := "5"

	tests := []struct {
		name          string
		cmd           CreateTemplateJawabanCommand
		mockJawaban   *mock.MockTemplateJawabanRepository
		mockPertanyan *mockTemplatePertanyaan.MockTemplatePertanyaanRepository
		expectedErr   string
	}{
		{
			name: "Success case",
			cmd: CreateTemplateJawabanCommand{
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				Jawaban:                "Sangat Baik",
				Nilai:                  &validNilaiStr,
				IsFreeText:             "0",
				SID:                    "sid-1",
				Resource:               "local",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				CreateFunc: func(ctx context.Context, templatejawaban *domaintemplatejawaban.TemplateJawaban) error {
					return nil
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   12,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "",
		},
		{
			name: "Invalid UuidTemplatePertanyaan",
			cmd: CreateTemplateJawabanCommand{
				UuidTemplatePertanyaan: "invalid-uuid",
				Jawaban:                "Sangat Baik",
				IsFreeText:             "0",
				SID:                    "sid-1",
				Resource:               "local",
			},
			mockJawaban:   &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{},
			expectedErr:   "TemplateJawaban.InvalidTemplatePertanyaan",
		},
		{
			name: "Pertanyaan not found",
			cmd: CreateTemplateJawabanCommand{
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				Jawaban:                "Sangat Baik",
				IsFreeText:             "0",
				SID:                    "sid-1",
				Resource:               "local",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return nil, errors.New("not found")
				},
			},
			expectedErr: "TemplateJawaban.NotFoundTemplatePertanyaan",
		},
		{
			name: "Invalid IsFreeText format",
			cmd: CreateTemplateJawabanCommand{
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				Jawaban:                "Sangat Baik",
				IsFreeText:             "not-a-number",
				SID:                    "sid-1",
				Resource:               "local",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   12,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "Must be a positive number",
		},
		{
			name: "Invalid owner (Resource is not local)",
			cmd: CreateTemplateJawabanCommand{
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				Jawaban:                "Sangat Baik",
				IsFreeText:             "0",
				SID:                    "sid-1",
				Resource:               "lpm",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   12,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "TemplateJawaban.InvalidOwner",
		},
		{
			name: "Repo Create failure",
			cmd: CreateTemplateJawabanCommand{
				UuidTemplatePertanyaan: validPertanyaanUUID.String(),
				Jawaban:                "Sangat Baik",
				Nilai:                  &validNilaiStr,
				IsFreeText:             "0",
				SID:                    "sid-1",
				Resource:               "local",
			},
			mockJawaban: &mock.MockTemplateJawabanRepository{
				CreateFunc: func(ctx context.Context, templatejawaban *domaintemplatejawaban.TemplateJawaban) error {
					return errors.New("database connection failed")
				},
			},
			mockPertanyan: &mockTemplatePertanyaan.MockTemplatePertanyaanRepository{
				GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domaintemplatepertanyaan.TemplatePertanyaan, error) {
					return &domaintemplatepertanyaan.TemplatePertanyaan{
						ID:   12,
						UUID: validPertanyaanUUID,
					}, nil
				},
			},
			expectedErr: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &CreateTemplateJawabanCommandHandler{
				Repo:           tt.mockJawaban,
				RepoPertanyaan: tt.mockPertanyan,
			}
			uuidStr, err := handler.Handle(context.Background(), tt.cmd)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, uuidStr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, uuidStr)
				_, parseErr := uuid.Parse(uuidStr)
				assert.NoError(t, parseErr)
			}
		})
	}
}
