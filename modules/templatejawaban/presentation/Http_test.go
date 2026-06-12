package presentation_test

import (
	"context"
	"encoding/json"
	"errors"
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
	"UnpakSiamida/modules/templatejawaban/domain"
	"UnpakSiamida/modules/templatejawaban/presentation"

	CreateTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/CreateTemplateJawaban"
	DeleteTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/DeleteTemplateJawaban"
	GetAllTemplateJawabans "UnpakSiamida/modules/templatejawaban/application/GetAllTemplateJawabans"
	GetTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/GetTemplateJawaban"
	RestoreTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/RestoreTemplateJawaban"
	SetupUuidTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/SetupUuidTemplateJawaban"
	UpdateTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/UpdateTemplateJawaban"
)

var (
	mockCreateTemplateJawabanFunc    func(ctx context.Context, cmd CreateTemplateJawaban.CreateTemplateJawabanCommand) (string, error)
	mockUpdateTemplateJawabanFunc    func(ctx context.Context, cmd UpdateTemplateJawaban.UpdateTemplateJawabanCommand) (string, error)
	mockDeleteTemplateJawabanFunc    func(ctx context.Context, cmd DeleteTemplateJawaban.DeleteTemplateJawabanCommand) (string, error)
	mockRestoreTemplateJawabanFunc   func(ctx context.Context, cmd RestoreTemplateJawaban.RestoreTemplateJawabanCommand) (string, error)
	mockGetTemplateJawabanFunc       func(ctx context.Context, q GetTemplateJawaban.GetTemplateJawabanByUuidQuery) (*domain.TemplateJawaban, error)
	mockGetAllTemplateJawabansFunc   func(ctx context.Context, q GetAllTemplateJawabans.GetAllTemplateJawabansQuery) (commondomain.Paged[domain.TemplateJawabanDefault], error)
	mockSetupUuidTemplateJawabanFunc func(ctx context.Context, cmd SetupUuidTemplateJawaban.SetupUuidTemplateJawabanCommand) (string, error)

	registerOnce sync.Once
)

type mockCreateTemplateJawabanHandler struct{}

func (h *mockCreateTemplateJawabanHandler) Handle(ctx context.Context, cmd CreateTemplateJawaban.CreateTemplateJawabanCommand) (string, error) {
	return mockCreateTemplateJawabanFunc(ctx, cmd)
}

type mockUpdateTemplateJawabanHandler struct{}

func (h *mockUpdateTemplateJawabanHandler) Handle(ctx context.Context, cmd UpdateTemplateJawaban.UpdateTemplateJawabanCommand) (string, error) {
	return mockUpdateTemplateJawabanFunc(ctx, cmd)
}

type mockDeleteTemplateJawabanHandler struct{}

func (h *mockDeleteTemplateJawabanHandler) Handle(ctx context.Context, cmd DeleteTemplateJawaban.DeleteTemplateJawabanCommand) (string, error) {
	return mockDeleteTemplateJawabanFunc(ctx, cmd)
}

type mockRestoreTemplateJawabanHandler struct{}

func (h *mockRestoreTemplateJawabanHandler) Handle(ctx context.Context, cmd RestoreTemplateJawaban.RestoreTemplateJawabanCommand) (string, error) {
	return mockRestoreTemplateJawabanFunc(ctx, cmd)
}

type mockGetTemplateJawabanHandler struct{}

func (h *mockGetTemplateJawabanHandler) Handle(ctx context.Context, q GetTemplateJawaban.GetTemplateJawabanByUuidQuery) (*domain.TemplateJawaban, error) {
	return mockGetTemplateJawabanFunc(ctx, q)
}

type mockGetAllTemplateJawabansHandler struct{}

func (h *mockGetAllTemplateJawabansHandler) Handle(ctx context.Context, q GetAllTemplateJawabans.GetAllTemplateJawabansQuery) (commondomain.Paged[domain.TemplateJawabanDefault], error) {
	return mockGetAllTemplateJawabansFunc(ctx, q)
}

type mockSetupUuidTemplateJawabanHandler struct{}

func (h *mockSetupUuidTemplateJawabanHandler) Handle(ctx context.Context, cmd SetupUuidTemplateJawaban.SetupUuidTemplateJawabanCommand) (string, error) {
	return mockSetupUuidTemplateJawabanFunc(ctx, cmd)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[CreateTemplateJawaban.CreateTemplateJawabanCommand, string](&mockCreateTemplateJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[UpdateTemplateJawaban.UpdateTemplateJawabanCommand, string](&mockUpdateTemplateJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteTemplateJawaban.DeleteTemplateJawabanCommand, string](&mockDeleteTemplateJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[RestoreTemplateJawaban.RestoreTemplateJawabanCommand, string](&mockRestoreTemplateJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[GetTemplateJawaban.GetTemplateJawabanByUuidQuery, *domain.TemplateJawaban](&mockGetTemplateJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllTemplateJawabans.GetAllTemplateJawabansQuery, commondomain.Paged[domain.TemplateJawabanDefault]](&mockGetAllTemplateJawabansHandler{})
		_ = mediatr.RegisterRequestHandler[SetupUuidTemplateJawaban.SetupUuidTemplateJawabanCommand, string](&mockSetupUuidTemplateJawabanHandler{})
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

func TestTemplateJawabanPresentation(t *testing.T) {
	setupMediatrMocks()

	app := fiber.New()
	presentation.ModuleTemplateJawaban(app)

	token := generateToken("user-123", "simak", "ctxFD49JawrQA")

	t.Run("Setup UUID TemplateJawaban success", func(t *testing.T) {
		mockSetupUuidTemplateJawabanFunc = func(ctx context.Context, cmd SetupUuidTemplateJawaban.SetupUuidTemplateJawabanCommand) (string, error) {
			return "setup complete", nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatejawaban/setupuuid", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "setup complete", res["message"])
	})

	t.Run("Create TemplateJawaban success", func(t *testing.T) {
		mockCreateTemplateJawabanFunc = func(ctx context.Context, cmd CreateTemplateJawaban.CreateTemplateJawabanCommand) (string, error) {
			assert.Equal(t, "pert-123", cmd.UuidTemplatePertanyaan)
			assert.Equal(t, "Sangat Baik", cmd.Jawaban)
			return "new-jawaban-uuid", nil
		}

		form := url.Values{}
		form.Add("template_pertanyaan", "pert-123")
		form.Add("jawaban", "Sangat Baik")

		req := httptest.NewRequest("POST", "/api/v2/templatejawaban", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-jawaban-uuid", res["uuid"])
	})

	t.Run("Update TemplateJawaban success", func(t *testing.T) {
		mockUpdateTemplateJawabanFunc = func(ctx context.Context, cmd UpdateTemplateJawaban.UpdateTemplateJawabanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "Cukup", cmd.Jawaban)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("jawaban", "Cukup")

		req := httptest.NewRequest("PUT", "/api/v2/templatejawaban/uuid-123", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete TemplateJawaban soft success", func(t *testing.T) {
		mockDeleteTemplateJawabanFunc = func(ctx context.Context, cmd DeleteTemplateJawaban.DeleteTemplateJawabanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "soft_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/templatejawaban/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete TemplateJawaban hard success", func(t *testing.T) {
		mockDeleteTemplateJawabanFunc = func(ctx context.Context, cmd DeleteTemplateJawaban.DeleteTemplateJawabanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "hard_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/templatejawaban/uuid-123/force", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Restore TemplateJawaban success", func(t *testing.T) {
		mockRestoreTemplateJawabanFunc = func(ctx context.Context, cmd RestoreTemplateJawaban.RestoreTemplateJawabanCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("PUT", "/api/v2/templatejawaban/uuid-123/restore", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Get TemplateJawaban by UUID success", func(t *testing.T) {
		jawabanVal := &domain.TemplateJawaban{
			Jawaban: "Baik",
		}
		mockGetTemplateJawabanFunc = func(ctx context.Context, q GetTemplateJawaban.GetTemplateJawabanByUuidQuery) (*domain.TemplateJawaban, error) {
			assert.Equal(t, "uuid-123", q.Uuid)
			return jawabanVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatejawaban/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.TemplateJawaban
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Baik", res.Jawaban)
	})

	t.Run("Get TemplateJawaban by UUID not found", func(t *testing.T) {
		mockGetTemplateJawabanFunc = func(ctx context.Context, q GetTemplateJawaban.GetTemplateJawabanByUuidQuery) (*domain.TemplateJawaban, error) {
			return nil, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatejawaban/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get TemplateJawaban by UUID error", func(t *testing.T) {
		mockGetTemplateJawabanFunc = func(ctx context.Context, q GetTemplateJawaban.GetTemplateJawabanByUuidQuery) (*domain.TemplateJawaban, error) {
			return nil, errors.New("database connection down")
		}

		req := httptest.NewRequest("GET", "/api/v2/templatejawaban/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Get All TemplateJawabans success", func(t *testing.T) {
		items := []domain.TemplateJawabanDefault{
			{Jawaban: "Baik"},
			{Jawaban: "Buruk"},
		}
		mockGetAllTemplateJawabansFunc = func(ctx context.Context, q GetAllTemplateJawabans.GetAllTemplateJawabansQuery) (commondomain.Paged[domain.TemplateJawabanDefault], error) {
			return commondomain.Paged[domain.TemplateJawabanDefault]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/templatejawabans?mode=all", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.TemplateJawabanDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 2)
		assert.Equal(t, "Baik", res[0].Jawaban)
	})
}
