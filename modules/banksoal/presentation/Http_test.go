package presentation_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/banksoal/domain"
	"UnpakSiamida/modules/banksoal/presentation"

	CopyBankSoal "UnpakSiamida/modules/banksoal/application/CopyBankSoal"
	CreateBankSoal "UnpakSiamida/modules/banksoal/application/CreateBankSoal"
	DeleteBankSoal "UnpakSiamida/modules/banksoal/application/DeleteBankSoal"
	DeleteTimeBankSoal "UnpakSiamida/modules/banksoal/application/DeleteTimeBankSoal"
	GetAllBankSoals "UnpakSiamida/modules/banksoal/application/GetAllBankSoals"
	GetBankSoalDefault "UnpakSiamida/modules/banksoal/application/GetBankSoalDefault"
	RestoreBankSoal "UnpakSiamida/modules/banksoal/application/RestoreBankSoal"
	ScheduleTimeBankSoal "UnpakSiamida/modules/banksoal/application/ScheduleTimeBankSoal"
	SetupUuidBankSoal "UnpakSiamida/modules/banksoal/application/SetupUuidBankSoal"
	StatusBankSoal "UnpakSiamida/modules/banksoal/application/StatusBankSoal"
	UpdateBankSoal "UnpakSiamida/modules/banksoal/application/UpdateBankSoal"
)

var (
	mockCreateBankSoalFunc       func(ctx context.Context, cmd CreateBankSoal.CreateBankSoalCommand) (string, error)
	mockChangeTimeBankSoalFunc   func(ctx context.Context, cmd ScheduleTimeBankSoal.ScheduleTimeBankSoalCommand) (string, error)
	mockUpdateBankSoalFunc       func(ctx context.Context, cmd UpdateBankSoal.UpdateBankSoalCommand) (string, error)
	mockStatusBankSoalFunc       func(ctx context.Context, cmd StatusBankSoal.StatusBankSoalCommand) (string, error)
	mockDeleteBankSoalFunc       func(ctx context.Context, cmd DeleteBankSoal.DeleteBankSoalCommand) (string, error)
	mockDeleteTimeBankSoalFunc   func(ctx context.Context, cmd DeleteTimeBankSoal.DeleteTimeBankSoalCommand) (string, error)
	mockDeleteTimeExtBankSoalFunc func(ctx context.Context, cmd DeleteTimeBankSoal.DeleteTimeExtBankSoalCommand) (string, error)
	mockRestoreBankSoalFunc      func(ctx context.Context, cmd RestoreBankSoal.RestoreBankSoalCommand) (string, error)
	mockCopyBankSoalFunc         func(ctx context.Context, cmd CopyBankSoal.CopyBankSoalCommand) (string, error)
	mockGetBankSoalDefaultFunc   func(ctx context.Context, q GetBankSoalDefault.GetBankSoalDefaultByUuidQuery) (*domain.BankSoalDefault, error)
	mockGetAllBankSoalsFunc      func(ctx context.Context, q GetAllBankSoals.GetAllBankSoalsQuery) (commondomain.Paged[domain.BankSoalDefault], error)
	mockSetupUuidBankSoalFunc    func(ctx context.Context, cmd SetupUuidBankSoal.SetupUuidBankSoalCommand) (string, error)

	registerOnce sync.Once
)

type mockCreateBankSoalHandler struct{}

func (h *mockCreateBankSoalHandler) Handle(ctx context.Context, cmd CreateBankSoal.CreateBankSoalCommand) (string, error) {
	return mockCreateBankSoalFunc(ctx, cmd)
}

type mockChangeTimeBankSoalHandler struct{}

func (h *mockChangeTimeBankSoalHandler) Handle(ctx context.Context, cmd ScheduleTimeBankSoal.ScheduleTimeBankSoalCommand) (string, error) {
	return mockChangeTimeBankSoalFunc(ctx, cmd)
}

type mockUpdateBankSoalHandler struct{}

func (h *mockUpdateBankSoalHandler) Handle(ctx context.Context, cmd UpdateBankSoal.UpdateBankSoalCommand) (string, error) {
	return mockUpdateBankSoalFunc(ctx, cmd)
}

type mockStatusBankSoalHandler struct{}

func (h *mockStatusBankSoalHandler) Handle(ctx context.Context, cmd StatusBankSoal.StatusBankSoalCommand) (string, error) {
	return mockStatusBankSoalFunc(ctx, cmd)
}

type mockDeleteBankSoalHandler struct{}

func (h *mockDeleteBankSoalHandler) Handle(ctx context.Context, cmd DeleteBankSoal.DeleteBankSoalCommand) (string, error) {
	return mockDeleteBankSoalFunc(ctx, cmd)
}

type mockDeleteTimeBankSoalHandler struct{}

func (h *mockDeleteTimeBankSoalHandler) Handle(ctx context.Context, cmd DeleteTimeBankSoal.DeleteTimeBankSoalCommand) (string, error) {
	return mockDeleteTimeBankSoalFunc(ctx, cmd)
}

type mockDeleteTimeExtBankSoalHandler struct{}

func (h *mockDeleteTimeExtBankSoalHandler) Handle(ctx context.Context, cmd DeleteTimeBankSoal.DeleteTimeExtBankSoalCommand) (string, error) {
	return mockDeleteTimeExtBankSoalFunc(ctx, cmd)
}

type mockRestoreBankSoalHandler struct{}

func (h *mockRestoreBankSoalHandler) Handle(ctx context.Context, cmd RestoreBankSoal.RestoreBankSoalCommand) (string, error) {
	return mockRestoreBankSoalFunc(ctx, cmd)
}

type mockCopyBankSoalHandler struct{}

func (h *mockCopyBankSoalHandler) Handle(ctx context.Context, cmd CopyBankSoal.CopyBankSoalCommand) (string, error) {
	return mockCopyBankSoalFunc(ctx, cmd)
}

type mockGetBankSoalDefaultHandler struct{}

func (h *mockGetBankSoalDefaultHandler) Handle(ctx context.Context, q GetBankSoalDefault.GetBankSoalDefaultByUuidQuery) (*domain.BankSoalDefault, error) {
	return mockGetBankSoalDefaultFunc(ctx, q)
}

type mockGetAllBankSoalsHandler struct{}

func (h *mockGetAllBankSoalsHandler) Handle(ctx context.Context, q GetAllBankSoals.GetAllBankSoalsQuery) (commondomain.Paged[domain.BankSoalDefault], error) {
	return mockGetAllBankSoalsFunc(ctx, q)
}

type mockSetupUuidBankSoalHandler struct{}

func (h *mockSetupUuidBankSoalHandler) Handle(ctx context.Context, cmd SetupUuidBankSoal.SetupUuidBankSoalCommand) (string, error) {
	return mockSetupUuidBankSoalFunc(ctx, cmd)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[CreateBankSoal.CreateBankSoalCommand, string](&mockCreateBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[ScheduleTimeBankSoal.ScheduleTimeBankSoalCommand, string](&mockChangeTimeBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[UpdateBankSoal.UpdateBankSoalCommand, string](&mockUpdateBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[StatusBankSoal.StatusBankSoalCommand, string](&mockStatusBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteBankSoal.DeleteBankSoalCommand, string](&mockDeleteBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteTimeBankSoal.DeleteTimeBankSoalCommand, string](&mockDeleteTimeBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteTimeBankSoal.DeleteTimeExtBankSoalCommand, string](&mockDeleteTimeExtBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[RestoreBankSoal.RestoreBankSoalCommand, string](&mockRestoreBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[CopyBankSoal.CopyBankSoalCommand, string](&mockCopyBankSoalHandler{})
		_ = mediatr.RegisterRequestHandler[GetBankSoalDefault.GetBankSoalDefaultByUuidQuery, *domain.BankSoalDefault](&mockGetBankSoalDefaultHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllBankSoals.GetAllBankSoalsQuery, commondomain.Paged[domain.BankSoalDefault]](&mockGetAllBankSoalsHandler{})
		_ = mediatr.RegisterRequestHandler[SetupUuidBankSoal.SetupUuidBankSoalCommand, string](&mockSetupUuidBankSoalHandler{})
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

func TestBankSoalPresentation(t *testing.T) {
	setupMediatrMocks()

	mockWhoamiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"ID": "dosen-123",
			"Level": "admin",
			"Resource": "simak",
			"CodeCtx": "ctxFD49JawrQA",
			"RefFakultas": "fak-123",
			"RefProdi": "prod-123"
		}`))
	}))
	defer mockWhoamiServer.Close()

	os.Setenv("WHOAMI_URL", mockWhoamiServer.URL)

	app := fiber.New()
	presentation.ModuleBankSoal(app)

	token := generateToken("user-123", "simak", "ctxFD49JawrQA")

	t.Run("Setup UUID BankSoal success", func(t *testing.T) {
		mockSetupUuidBankSoalFunc = func(ctx context.Context, cmd SetupUuidBankSoal.SetupUuidBankSoalCommand) (string, error) {
			return "setup complete", nil
		}

		req := httptest.NewRequest("GET", "/api/v2/banksoal/setupuuid", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "setup complete", res["message"])
	})

	t.Run("Create BankSoal success", func(t *testing.T) {
		mockCreateBankSoalFunc = func(ctx context.Context, cmd CreateBankSoal.CreateBankSoalCommand) (string, error) {
			assert.Equal(t, "UAS Math", cmd.Judul)
			return "new-banksoal-uuid", nil
		}

		form := url.Values{}
		form.Add("judul", "UAS Math")
		form.Add("content", "Some content")

		req := httptest.NewRequest("POST", "/api/v2/banksoal", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-banksoal-uuid", res["uuid"])
	})

	t.Run("Change Time BankSoal success", func(t *testing.T) {
		mockChangeTimeBankSoalFunc = func(ctx context.Context, cmd ScheduleTimeBankSoal.ScheduleTimeBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.UuidBankSoal)
			assert.Equal(t, "2026-06-12", cmd.TanggalMulai)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("tanggal_mulai", "2026-06-12")
		form.Add("tanggal_akhir", "2026-06-13")

		req := httptest.NewRequest("PUT", "/api/v2/banksoal/uuid-123/schedule", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Update BankSoal success", func(t *testing.T) {
		mockUpdateBankSoalFunc = func(ctx context.Context, cmd UpdateBankSoal.UpdateBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "New Judul", cmd.Judul)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("judul", "New Judul")

		req := httptest.NewRequest("PUT", "/api/v2/banksoal/uuid-123", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Status BankSoal success", func(t *testing.T) {
		mockStatusBankSoalFunc = func(ctx context.Context, cmd StatusBankSoal.StatusBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "active", cmd.Status)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("status", "active")

		req := httptest.NewRequest("PUT", "/api/v2/banksoal/uuid-123/status", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete BankSoal soft success", func(t *testing.T) {
		mockDeleteBankSoalFunc = func(ctx context.Context, cmd DeleteBankSoal.DeleteBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "soft_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/banksoal/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete BankSoal hard success", func(t *testing.T) {
		mockDeleteBankSoalFunc = func(ctx context.Context, cmd DeleteBankSoal.DeleteBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "hard_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/banksoal/uuid-123/force", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete Time BankSoal success", func(t *testing.T) {
		mockDeleteTimeBankSoalFunc = func(ctx context.Context, cmd DeleteTimeBankSoal.DeleteTimeBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/banksoal/uuid-123/time", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete Time Ext BankSoal success", func(t *testing.T) {
		mockDeleteTimeExtBankSoalFunc = func(ctx context.Context, cmd DeleteTimeBankSoal.DeleteTimeExtBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "soal-abc", cmd.UuidBankSoal)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("banksoal", "soal-abc")

		req := httptest.NewRequest("DELETE", "/api/v2/banksoal/uuid-123/timeext", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Restore BankSoal success", func(t *testing.T) {
		mockRestoreBankSoalFunc = func(ctx context.Context, cmd RestoreBankSoal.RestoreBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("PUT", "/api/v2/banksoal/uuid-123/restore", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Copy BankSoal success", func(t *testing.T) {
		mockCopyBankSoalFunc = func(ctx context.Context, cmd CopyBankSoal.CopyBankSoalCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "new-copy-uuid", nil
		}

		req := httptest.NewRequest("POST", "/api/v2/banksoal/uuid-123/copy", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-copy-uuid", res["uuid"])
	})

	t.Run("Get BankSoal by UUID success", func(t *testing.T) {
		banksoalVal := &domain.BankSoalDefault{
			Judul: "Final Exam",
		}
		mockGetBankSoalDefaultFunc = func(ctx context.Context, q GetBankSoalDefault.GetBankSoalDefaultByUuidQuery) (*domain.BankSoalDefault, error) {
			assert.Equal(t, "uuid-123", q.Uuid)
			return banksoalVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/banksoal/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.BankSoalDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Final Exam", res.Judul)
	})

	t.Run("Get BankSoal by UUID not found", func(t *testing.T) {
		mockGetBankSoalDefaultFunc = func(ctx context.Context, q GetBankSoalDefault.GetBankSoalDefaultByUuidQuery) (*domain.BankSoalDefault, error) {
			return nil, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/banksoal/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get All BankSoals success", func(t *testing.T) {
		items := []domain.BankSoalDefault{
			{Judul: "Soal 1"},
			{Judul: "Soal 2"},
		}
		mockGetAllBankSoalsFunc = func(ctx context.Context, q GetAllBankSoals.GetAllBankSoalsQuery) (commondomain.Paged[domain.BankSoalDefault], error) {
			assert.Equal(t, "prod-123", *q.TargetProdi)
			return commondomain.Paged[domain.BankSoalDefault]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/banksoals?mode=all", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.BankSoalDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 2)
		assert.Equal(t, "Soal 1", res[0].Judul)
	})
}
