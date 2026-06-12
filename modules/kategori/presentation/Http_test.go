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
	"UnpakSiamida/modules/kategori/domain"
	"UnpakSiamida/modules/kategori/presentation"

	CopyKategori "UnpakSiamida/modules/kategori/application/CopyKategori"
	CreateKategori "UnpakSiamida/modules/kategori/application/CreateKategori"
	DeleteKategori "UnpakSiamida/modules/kategori/application/DeleteKategori"
	GetAllKategoris "UnpakSiamida/modules/kategori/application/GetAllKategoris"
	GetKategori "UnpakSiamida/modules/kategori/application/GetKategori"
	RestoreKategori "UnpakSiamida/modules/kategori/application/RestoreKategori"
	SetupUuidKategori "UnpakSiamida/modules/kategori/application/SetupUuidKategori"
	UpdateKategori "UnpakSiamida/modules/kategori/application/UpdateKategori"
)

var (
	mockCreateKategoriFunc      func(ctx context.Context, cmd CreateKategori.CreateKategoriCommand) (string, error)
	mockUpdateKategoriOrderFunc func(ctx context.Context, cmd UpdateKategori.UpdateKategoriOrderCommand) (string, error)
	mockUpdateKategoriFunc      func(ctx context.Context, cmd UpdateKategori.UpdateKategoriCommand) (string, error)
	mockDeleteKategoriFunc      func(ctx context.Context, cmd DeleteKategori.DeleteKategoriCommand) (string, error)
	mockRestoreKategoriFunc     func(ctx context.Context, cmd RestoreKategori.RestoreKategoriCommand) (string, error)
	mockCopyKategoriFunc        func(ctx context.Context, cmd CopyKategori.CopyKategoriCommand) (string, error)
	mockGetKategoriByUuidFunc   func(ctx context.Context, q GetKategori.GetKategoriByUuidQuery) (*domain.Kategori, error)
	mockGetAllKategorisFunc     func(ctx context.Context, q GetAllKategoris.GetAllKategorisQuery) (commondomain.Paged[domain.KategoriDefault], error)
	mockSetupUuidKategoriFunc   func(ctx context.Context, cmd SetupUuidKategori.SetupUuidKategoriCommand) (string, error)

	registerOnce sync.Once
)

type mockCreateKategoriHandler struct{}

func (h *mockCreateKategoriHandler) Handle(ctx context.Context, cmd CreateKategori.CreateKategoriCommand) (string, error) {
	return mockCreateKategoriFunc(ctx, cmd)
}

type mockUpdateKategoriOrderHandler struct{}

func (h *mockUpdateKategoriOrderHandler) Handle(ctx context.Context, cmd UpdateKategori.UpdateKategoriOrderCommand) (string, error) {
	return mockUpdateKategoriOrderFunc(ctx, cmd)
}

type mockUpdateKategoriHandler struct{}

func (h *mockUpdateKategoriHandler) Handle(ctx context.Context, cmd UpdateKategori.UpdateKategoriCommand) (string, error) {
	return mockUpdateKategoriFunc(ctx, cmd)
}

type mockDeleteKategoriHandler struct{}

func (h *mockDeleteKategoriHandler) Handle(ctx context.Context, cmd DeleteKategori.DeleteKategoriCommand) (string, error) {
	return mockDeleteKategoriFunc(ctx, cmd)
}

type mockRestoreKategoriHandler struct{}

func (h *mockRestoreKategoriHandler) Handle(ctx context.Context, cmd RestoreKategori.RestoreKategoriCommand) (string, error) {
	return mockRestoreKategoriFunc(ctx, cmd)
}

type mockCopyKategoriHandler struct{}

func (h *mockCopyKategoriHandler) Handle(ctx context.Context, cmd CopyKategori.CopyKategoriCommand) (string, error) {
	return mockCopyKategoriFunc(ctx, cmd)
}

type mockGetKategoriByUuidHandler struct{}

func (h *mockGetKategoriByUuidHandler) Handle(ctx context.Context, q GetKategori.GetKategoriByUuidQuery) (*domain.Kategori, error) {
	return mockGetKategoriByUuidFunc(ctx, q)
}

type mockGetAllKategorisHandler struct{}

func (h *mockGetAllKategorisHandler) Handle(ctx context.Context, q GetAllKategoris.GetAllKategorisQuery) (commondomain.Paged[domain.KategoriDefault], error) {
	return mockGetAllKategorisFunc(ctx, q)
}

type mockSetupUuidKategoriHandler struct{}

func (h *mockSetupUuidKategoriHandler) Handle(ctx context.Context, cmd SetupUuidKategori.SetupUuidKategoriCommand) (string, error) {
	return mockSetupUuidKategoriFunc(ctx, cmd)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[CreateKategori.CreateKategoriCommand, string](&mockCreateKategoriHandler{})
		_ = mediatr.RegisterRequestHandler[UpdateKategori.UpdateKategoriOrderCommand, string](&mockUpdateKategoriOrderHandler{})
		_ = mediatr.RegisterRequestHandler[UpdateKategori.UpdateKategoriCommand, string](&mockUpdateKategoriHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteKategori.DeleteKategoriCommand, string](&mockDeleteKategoriHandler{})
		_ = mediatr.RegisterRequestHandler[RestoreKategori.RestoreKategoriCommand, string](&mockRestoreKategoriHandler{})
		_ = mediatr.RegisterRequestHandler[CopyKategori.CopyKategoriCommand, string](&mockCopyKategoriHandler{})
		_ = mediatr.RegisterRequestHandler[GetKategori.GetKategoriByUuidQuery, *domain.Kategori](&mockGetKategoriByUuidHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllKategoris.GetAllKategorisQuery, commondomain.Paged[domain.KategoriDefault]](&mockGetAllKategorisHandler{})
		_ = mediatr.RegisterRequestHandler[SetupUuidKategori.SetupUuidKategoriCommand, string](&mockSetupUuidKategoriHandler{})
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

func TestKategoriPresentation(t *testing.T) {
	setupMediatrMocks()

	app := fiber.New()
	presentation.ModuleKategori(app)

	token := generateToken("user-123", "simak", "dosen")

	t.Run("Setup UUID Kategori success", func(t *testing.T) {
		mockSetupUuidKategoriFunc = func(ctx context.Context, cmd SetupUuidKategori.SetupUuidKategoriCommand) (string, error) {
			return "setup complete", nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kategori/setupuuid", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "setup complete", res["message"])
	})

	t.Run("Create Kategori success", func(t *testing.T) {
		mockCreateKategoriFunc = func(ctx context.Context, cmd CreateKategori.CreateKategoriCommand) (string, error) {
			assert.Equal(t, "Test Kategori", cmd.NamaKategori)
			return "new-kategori-uuid", nil
		}

		form := url.Values{}
		form.Add("nama_kategori", "Test Kategori")

		req := httptest.NewRequest("POST", "/api/v2/kategori", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-kategori-uuid", res["uuid"])
	})

	t.Run("Update Kategori Order success", func(t *testing.T) {
		mockUpdateKategoriOrderFunc = func(ctx context.Context, cmd UpdateKategori.UpdateKategoriOrderCommand) (string, error) {
			assert.Equal(t, "order-payload", cmd.Payload)
			return "updated-uuid", nil
		}

		form := url.Values{}
		form.Add("payload", "order-payload")

		req := httptest.NewRequest("PUT", "/api/v2/kategori", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "updated-uuid", res["uuid"])
	})

	t.Run("Update Kategori success", func(t *testing.T) {
		mockUpdateKategoriFunc = func(ctx context.Context, cmd UpdateKategori.UpdateKategoriCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "Updated Name", cmd.NamaKategori)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("nama_kategori", "Updated Name")

		req := httptest.NewRequest("PUT", "/api/v2/kategori/uuid-123", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete Kategori soft success", func(t *testing.T) {
		mockDeleteKategoriFunc = func(ctx context.Context, cmd DeleteKategori.DeleteKategoriCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "soft_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/kategori/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete Kategori hard success", func(t *testing.T) {
		mockDeleteKategoriFunc = func(ctx context.Context, cmd DeleteKategori.DeleteKategoriCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "hard_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/kategori/uuid-123/force", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Restore Kategori success", func(t *testing.T) {
		mockRestoreKategoriFunc = func(ctx context.Context, cmd RestoreKategori.RestoreKategoriCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("PUT", "/api/v2/kategori/uuid-123/restore", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Copy Kategori success", func(t *testing.T) {
		mockCopyKategoriFunc = func(ctx context.Context, cmd CopyKategori.CopyKategoriCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "new-copy-uuid", nil
		}

		req := httptest.NewRequest("POST", "/api/v2/kategori/uuid-123/copy", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-copy-uuid", res["uuid"])
	})

	t.Run("Get Kategori by UUID success", func(t *testing.T) {
		kategoriVal := &domain.Kategori{
			NamaKategori: "Math",
		}
		mockGetKategoriByUuidFunc = func(ctx context.Context, q GetKategori.GetKategoriByUuidQuery) (*domain.Kategori, error) {
			assert.Equal(t, "uuid-123", q.Uuid)
			return kategoriVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kategori/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.Kategori
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Math", res.NamaKategori)
	})

	t.Run("Get Kategori by UUID not found", func(t *testing.T) {
		mockGetKategoriByUuidFunc = func(ctx context.Context, q GetKategori.GetKategoriByUuidQuery) (*domain.Kategori, error) {
			return nil, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kategori/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get Kategori by UUID error", func(t *testing.T) {
		mockGetKategoriByUuidFunc = func(ctx context.Context, q GetKategori.GetKategoriByUuidQuery) (*domain.Kategori, error) {
			return nil, errors.New("database connection down")
		}

		req := httptest.NewRequest("GET", "/api/v2/kategori/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Get All Kategoris success", func(t *testing.T) {
		items := []domain.KategoriDefault{
			{NamaKategori: "Cat 1"},
			{NamaKategori: "Cat 2"},
		}
		mockGetAllKategorisFunc = func(ctx context.Context, q GetAllKategoris.GetAllKategorisQuery) (commondomain.Paged[domain.KategoriDefault], error) {
			return commondomain.Paged[domain.KategoriDefault]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kategoris?mode=all", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.KategoriDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 2)
		assert.Equal(t, "Cat 1", res[0].NamaKategori)
	})
}
