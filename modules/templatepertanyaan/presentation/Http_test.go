package presentation_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/templatepertanyaan/domain"
	"UnpakSiamida/modules/templatepertanyaan/presentation"

	CopyTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/CopyTemplatePertanyaan"
	CreateTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/CreateTemplatePertanyaan"
	DeleteTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/DeleteTemplatePertanyaan"
	GetAllTemplatePertanyaans "UnpakSiamida/modules/templatepertanyaan/application/GetAllTemplatePertanyaans"
	GetTemplatePertanyaanDefault "UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaanDefault"
	GetTemplate "UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaanWithAnswareDefault"
	RestoreTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/RestoreTemplatePertanyaan"
	SetupUuidTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/SetupUuidTemplatePertanyaan"
	StatusTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/StatusTemplatePertanyaan"
	UpdateTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/UpdateTemplatePertanyaan"
)

var (
	mockCreateTemplatePertanyaanFunc             func(ctx context.Context, cmd CreateTemplatePertanyaan.CreateTemplatePertanyaanCommand) (string, error)
	mockUpdateTemplatePertanyaanFunc             func(ctx context.Context, cmd UpdateTemplatePertanyaan.UpdateTemplatePertanyaanCommand) (string, error)
	mockStatusTemplatePertanyaanFunc             func(ctx context.Context, cmd StatusTemplatePertanyaan.StatusTemplatePertanyaanCommand) (string, error)
	mockDeleteTemplatePertanyaanFunc             func(ctx context.Context, cmd DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand) (string, error)
	mockRestoreTemplatePertanyaanFunc            func(ctx context.Context, cmd RestoreTemplatePertanyaan.RestoreTemplatePertanyaanCommand) (string, error)
	mockCopyTemplatePertanyaanFunc               func(ctx context.Context, cmd CopyTemplatePertanyaan.CopyTemplatePertanyaanCommand) (string, error)
	mockGetTemplatePertanyaanDefaultFunc         func(ctx context.Context, q GetTemplatePertanyaanDefault.GetTemplatePertanyaanDefaultByUuidQuery) (*domain.TemplatePertanyaanDefault, error)
	mockGetTemplatePertanyaanWithAnswareFunc     func(ctx context.Context, q GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery) (*domain.TemplatePertanyaanWithAnswareDefault, error)
	mockGetTemplatePertanyaanByBankSoalFunc      func(ctx context.Context, q GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery) (commondomain.Paged[domain.TemplatePertanyaanWithAnswareDefault], error)
	mockGetAllTemplatePertanyaansFunc            func(ctx context.Context, q GetAllTemplatePertanyaans.GetAllTemplatePertanyaansQuery) (commondomain.Paged[domain.TemplatePertanyaanDefault], error)
	mockSetupUuidTemplatePertanyaanFunc          func(ctx context.Context, cmd SetupUuidTemplatePertanyaan.SetupUuidTemplatePertanyaanCommand) (string, error)

	registerOnce sync.Once
)

type mockCreateTemplatePertanyaanHandler struct{}

func (h *mockCreateTemplatePertanyaanHandler) Handle(ctx context.Context, cmd CreateTemplatePertanyaan.CreateTemplatePertanyaanCommand) (string, error) {
	return mockCreateTemplatePertanyaanFunc(ctx, cmd)
}

type mockUpdateTemplatePertanyaanHandler struct{}

func (h *mockUpdateTemplatePertanyaanHandler) Handle(ctx context.Context, cmd UpdateTemplatePertanyaan.UpdateTemplatePertanyaanCommand) (string, error) {
	return mockUpdateTemplatePertanyaanFunc(ctx, cmd)
}

type mockStatusTemplatePertanyaanHandler struct{}

func (h *mockStatusTemplatePertanyaanHandler) Handle(ctx context.Context, cmd StatusTemplatePertanyaan.StatusTemplatePertanyaanCommand) (string, error) {
	return mockStatusTemplatePertanyaanFunc(ctx, cmd)
}

type mockDeleteTemplatePertanyaanHandler struct{}

func (h *mockDeleteTemplatePertanyaanHandler) Handle(ctx context.Context, cmd DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand) (string, error) {
	return mockDeleteTemplatePertanyaanFunc(ctx, cmd)
}

type mockRestoreTemplatePertanyaanHandler struct{}

func (h *mockRestoreTemplatePertanyaanHandler) Handle(ctx context.Context, cmd RestoreTemplatePertanyaan.RestoreTemplatePertanyaanCommand) (string, error) {
	return mockRestoreTemplatePertanyaanFunc(ctx, cmd)
}

type mockCopyTemplatePertanyaanHandler struct{}

func (h *mockCopyTemplatePertanyaanHandler) Handle(ctx context.Context, cmd CopyTemplatePertanyaan.CopyTemplatePertanyaanCommand) (string, error) {
	return mockCopyTemplatePertanyaanFunc(ctx, cmd)
}

type mockGetTemplatePertanyaanDefaultHandler struct{}

func (h *mockGetTemplatePertanyaanDefaultHandler) Handle(ctx context.Context, q GetTemplatePertanyaanDefault.GetTemplatePertanyaanDefaultByUuidQuery) (*domain.TemplatePertanyaanDefault, error) {
	return mockGetTemplatePertanyaanDefaultFunc(ctx, q)
}

type mockGetTemplatePertanyaanWithAnswareHandler struct{}

func (h *mockGetTemplatePertanyaanWithAnswareHandler) Handle(ctx context.Context, q GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery) (*domain.TemplatePertanyaanWithAnswareDefault, error) {
	return mockGetTemplatePertanyaanWithAnswareFunc(ctx, q)
}

type mockGetTemplatePertanyaanByBankSoalHandler struct{}

func (h *mockGetTemplatePertanyaanByBankSoalHandler) Handle(ctx context.Context, q GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery) (commondomain.Paged[domain.TemplatePertanyaanWithAnswareDefault], error) {
	return mockGetTemplatePertanyaanByBankSoalFunc(ctx, q)
}

type mockGetAllTemplatePertanyaansHandler struct{}

func (h *mockGetAllTemplatePertanyaansHandler) Handle(ctx context.Context, q GetAllTemplatePertanyaans.GetAllTemplatePertanyaansQuery) (commondomain.Paged[domain.TemplatePertanyaanDefault], error) {
	return mockGetAllTemplatePertanyaansFunc(ctx, q)
}

type mockSetupUuidTemplatePertanyaanHandler struct{}

func (h *mockSetupUuidTemplatePertanyaanHandler) Handle(ctx context.Context, cmd SetupUuidTemplatePertanyaan.SetupUuidTemplatePertanyaanCommand) (string, error) {
	return mockSetupUuidTemplatePertanyaanFunc(ctx, cmd)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[CreateTemplatePertanyaan.CreateTemplatePertanyaanCommand, string](&mockCreateTemplatePertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[UpdateTemplatePertanyaan.UpdateTemplatePertanyaanCommand, string](&mockUpdateTemplatePertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[StatusTemplatePertanyaan.StatusTemplatePertanyaanCommand, string](&mockStatusTemplatePertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand, string](&mockDeleteTemplatePertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[RestoreTemplatePertanyaan.RestoreTemplatePertanyaanCommand, string](&mockRestoreTemplatePertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[CopyTemplatePertanyaan.CopyTemplatePertanyaanCommand, string](&mockCopyTemplatePertanyaanHandler{})
		_ = mediatr.RegisterRequestHandler[GetTemplatePertanyaanDefault.GetTemplatePertanyaanDefaultByUuidQuery, *domain.TemplatePertanyaanDefault](&mockGetTemplatePertanyaanDefaultHandler{})
		_ = mediatr.RegisterRequestHandler[GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery, *domain.TemplatePertanyaanWithAnswareDefault](&mockGetTemplatePertanyaanWithAnswareHandler{})
		_ = mediatr.RegisterRequestHandler[GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery, commondomain.Paged[domain.TemplatePertanyaanWithAnswareDefault]](&mockGetTemplatePertanyaanByBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllTemplatePertanyaans.GetAllTemplatePertanyaansQuery, commondomain.Paged[domain.TemplatePertanyaanDefault]](&mockGetAllTemplatePertanyaansHandler{})
		_ = mediatr.RegisterRequestHandler[SetupUuidTemplatePertanyaan.SetupUuidTemplatePertanyaanCommand, string](&mockSetupUuidTemplatePertanyaanHandler{})
	})
}

func generateToken(sid, resource, codectx string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":      sid,
		"resource": resource,
		"codectx":  codectx,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}

func TestTemplatePertanyaanPresentation(t *testing.T) {
	setupMediatrMocks()

	app := fiber.New()
	presentation.ModuleTemplatePertanyaan(app)

	token := generateToken("user-123", "simak", "ctxFD49JawrQA")

	t.Run("Setup UUID TemplatePertanyaan success", func(t *testing.T) {
		mockSetupUuidTemplatePertanyaanFunc = func(ctx context.Context, cmd SetupUuidTemplatePertanyaan.SetupUuidTemplatePertanyaanCommand) (string, error) {
			return "setup complete", nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatepertanyaan/setupuuid", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "setup complete", res["message"])
	})

	t.Run("Create TemplatePertanyaan success", func(t *testing.T) {
		mockCreateTemplatePertanyaanFunc = func(ctx context.Context, cmd CreateTemplatePertanyaan.CreateTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "soal-123", cmd.UuidBankSoal)
			assert.Equal(t, "Apakah dosen mengajar dengan baik?", cmd.Pertanyaan)
			return "new-pertanyaan-uuid", nil
		}

		form := url.Values{}
		form.Add("bank_soal", "soal-123")
		form.Add("pertanyaan", "Apakah dosen mengajar dengan baik?")
		form.Add("jenis_pilihan", "pilihan_ganda")
		form.Add("bobot", "5")
		form.Add("required", "1")

		req := httptest.NewRequest("POST", "/api/v2/templatepertanyaan", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-pertanyaan-uuid", res["uuid"])
	})

	t.Run("Update TemplatePertanyaan success", func(t *testing.T) {
		mockUpdateTemplatePertanyaanFunc = func(ctx context.Context, cmd UpdateTemplatePertanyaan.UpdateTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "Updated Pertanyaan", cmd.Pertanyaan)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("pertanyaan", "Updated Pertanyaan")

		req := httptest.NewRequest("PUT", "/api/v2/templatepertanyaan/uuid-123", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Status TemplatePertanyaan success", func(t *testing.T) {
		mockStatusTemplatePertanyaanFunc = func(ctx context.Context, cmd StatusTemplatePertanyaan.StatusTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "draft", cmd.Status)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("status", "draft")

		req := httptest.NewRequest("PUT", "/api/v2/templatepertanyaan/uuid-123/status", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete TemplatePertanyaan soft success", func(t *testing.T) {
		mockDeleteTemplatePertanyaanFunc = func(ctx context.Context, cmd DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "soft_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/templatepertanyaan/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete TemplatePertanyaan hard success", func(t *testing.T) {
		mockDeleteTemplatePertanyaanFunc = func(ctx context.Context, cmd DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "hard_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/templatepertanyaan/uuid-123/force", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Restore TemplatePertanyaan success", func(t *testing.T) {
		mockRestoreTemplatePertanyaanFunc = func(ctx context.Context, cmd RestoreTemplatePertanyaan.RestoreTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("PUT", "/api/v2/templatepertanyaan/uuid-123/restore", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Copy TemplatePertanyaan success", func(t *testing.T) {
		mockCopyTemplatePertanyaanFunc = func(ctx context.Context, cmd CopyTemplatePertanyaan.CopyTemplatePertanyaanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "new-copy-uuid", nil
		}

		req := httptest.NewRequest("POST", "/api/v2/templatepertanyaan/uuid-123/copy", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-copy-uuid", res["uuid"])
	})

	t.Run("Get TemplatePertanyaan by UUID success", func(t *testing.T) {
		pertanyaanVal := &domain.TemplatePertanyaanDefault{
			Pertanyaan: "Apakah mengajar?",
		}
		mockGetTemplatePertanyaanDefaultFunc = func(ctx context.Context, q GetTemplatePertanyaanDefault.GetTemplatePertanyaanDefaultByUuidQuery) (*domain.TemplatePertanyaanDefault, error) {
			assert.Equal(t, "uuid-123", q.Uuid)
			return pertanyaanVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatepertanyaan/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.TemplatePertanyaanDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Apakah mengajar?", res.Pertanyaan)
	})

	t.Run("Get Template by UUID success", func(t *testing.T) {
		withAnswareVal := &domain.TemplatePertanyaanWithAnswareDefault{
			Pertanyaan: "Dengan Jawaban",
		}
		mockGetTemplatePertanyaanWithAnswareFunc = func(ctx context.Context, q GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByUuidQuery) (*domain.TemplatePertanyaanWithAnswareDefault, error) {
			assert.Equal(t, "uuid-123", q.Uuid)
			return withAnswareVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatepertanyaan/uuid-123/template", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.TemplatePertanyaanWithAnswareDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Dengan Jawaban", res.Pertanyaan)
	})

	t.Run("Get Template Pertanyaan by BankSoal success", func(t *testing.T) {
		items := []domain.TemplatePertanyaanWithAnswareDefault{
			{Pertanyaan: "Soal 1"},
		}
		mockGetTemplatePertanyaanByBankSoalFunc = func(ctx context.Context, q GetTemplate.GetTemplatePertanyaanWithAnswareDefaultByBankSoalQuery) (commondomain.Paged[domain.TemplatePertanyaanWithAnswareDefault], error) {
			assert.Equal(t, "soal-123", q.UuidBankSoal)
			return commondomain.Paged[domain.TemplatePertanyaanWithAnswareDefault]{
				Data:  items,
				Total: 1,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatepertanyaan/soal-123/banksoal", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get All TemplatePertanyaans success", func(t *testing.T) {
		items := []domain.TemplatePertanyaanDefault{
			{Pertanyaan: "Pert 1"},
			{Pertanyaan: "Pert 2"},
		}
		mockGetAllTemplatePertanyaansFunc = func(ctx context.Context, q GetAllTemplatePertanyaans.GetAllTemplatePertanyaansQuery) (commondomain.Paged[domain.TemplatePertanyaanDefault], error) {
			return commondomain.Paged[domain.TemplatePertanyaanDefault]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatepertanyaans?mode=all", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.TemplatePertanyaanDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 2)
		assert.Equal(t, "Pert 1", res[0].Pertanyaan)
	})
}
