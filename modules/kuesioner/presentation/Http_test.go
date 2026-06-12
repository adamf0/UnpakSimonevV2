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
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	commondomain "UnpakSiamida/common/domain"
	BankSoaldomain "UnpakSiamida/modules/banksoal/domain"
	"UnpakSiamida/modules/kuesioner/domain"
	"UnpakSiamida/modules/kuesioner/presentation"

	ActiveKuesioners "UnpakSiamida/modules/kuesioner/application/ActiveKuesioner"
	ActiveKuesionersSingle "UnpakSiamida/modules/kuesioner/application/ActiveKuesionerSingle"
	CreateKuesioner "UnpakSiamida/modules/kuesioner/application/CreateKuesioner"
	DeleteKuesioner "UnpakSiamida/modules/kuesioner/application/DeleteKuesioner"
	GetAllKuesionerReport "UnpakSiamida/modules/kuesioner/application/GetAllKuesionerReport"
	GetAllKuesioners "UnpakSiamida/modules/kuesioner/application/GetAllKuesioners"
	GetKuesioner "UnpakSiamida/modules/kuesioner/application/GetKuesioner"
	GetKuesionerJawaban "UnpakSiamida/modules/kuesioner/application/GetKuesionerJawaban"
	SaveKuesionerJawaban "UnpakSiamida/modules/kuesioner/application/SaveKuesionerJawaban"
	SetupUuidKuesioner "UnpakSiamida/modules/kuesioner/application/SetupUuidKuesioner"
)

var (
	mockCreateKuesionerFunc          func(ctx context.Context, cmd CreateKuesioner.CreateKuesionerCommand) (string, error)
	mockSaveKuesionerJawabanFunc     func(ctx context.Context, cmd SaveKuesionerJawaban.SaveKuesionerJawabanCommand) (string, error)
	mockGetKuesionerJawabanFunc      func(ctx context.Context, q GetKuesionerJawaban.GetKuesionerJawabanByUuidQuery) ([]domain.KuesionerJawabanDefault, error)
	mockDeleteKuesionerFunc          func(ctx context.Context, cmd DeleteKuesioner.DeleteKuesionerCommand) (string, error)
	mockGetKuesionerFunc             func(ctx context.Context, q GetKuesioner.GetKuesionerByUuidQuery) (*domain.Kuesioner, error)
	mockGetAllKuesionersFunc         func(ctx context.Context, q GetAllKuesioners.GetAllKuesionersQuery) (commondomain.Paged[domain.KuesionerDefault], error)
	mockGetAllKuesionersReportFunc   func(ctx context.Context, q GetAllKuesionerReport.GetAllKuesionersReportQuery) ([]domain.KuesionerResult, error)
	mockActiveKuesionersFunc         func(ctx context.Context, q ActiveKuesioners.ActiveKuesionerQuery) ([]BankSoaldomain.BankSoalDefault, error)
	mockActiveKuesionersSingleFunc   func(ctx context.Context, q ActiveKuesionersSingle.ActiveKuesionerSingleQuery) (*BankSoaldomain.BankSoalDefault, error)
	mockSetupUuidKuesionerFunc       func(ctx context.Context, cmd SetupUuidKuesioner.SetupUuidKuesionerCommand) (string, error)

	registerOnce sync.Once
)

type mockCreateKuesionerHandler struct{}

func (h *mockCreateKuesionerHandler) Handle(ctx context.Context, cmd CreateKuesioner.CreateKuesionerCommand) (string, error) {
	return mockCreateKuesionerFunc(ctx, cmd)
}

type mockSaveKuesionerJawabanHandler struct{}

func (h *mockSaveKuesionerJawabanHandler) Handle(ctx context.Context, cmd SaveKuesionerJawaban.SaveKuesionerJawabanCommand) (string, error) {
	return mockSaveKuesionerJawabanFunc(ctx, cmd)
}

type mockGetKuesionerJawabanHandler struct{}

func (h *mockGetKuesionerJawabanHandler) Handle(ctx context.Context, q GetKuesionerJawaban.GetKuesionerJawabanByUuidQuery) ([]domain.KuesionerJawabanDefault, error) {
	return mockGetKuesionerJawabanFunc(ctx, q)
}

type mockDeleteKuesionerHandler struct{}

func (h *mockDeleteKuesionerHandler) Handle(ctx context.Context, cmd DeleteKuesioner.DeleteKuesionerCommand) (string, error) {
	return mockDeleteKuesionerFunc(ctx, cmd)
}

type mockGetKuesionerHandler struct{}

func (h *mockGetKuesionerHandler) Handle(ctx context.Context, q GetKuesioner.GetKuesionerByUuidQuery) (*domain.Kuesioner, error) {
	return mockGetKuesionerFunc(ctx, q)
}

type mockGetAllKuesionersHandler struct{}

func (h *mockGetAllKuesionersHandler) Handle(ctx context.Context, q GetAllKuesioners.GetAllKuesionersQuery) (commondomain.Paged[domain.KuesionerDefault], error) {
	return mockGetAllKuesionersFunc(ctx, q)
}

type mockGetAllKuesionersReportHandler struct{}

func (h *mockGetAllKuesionersReportHandler) Handle(ctx context.Context, q GetAllKuesionerReport.GetAllKuesionersReportQuery) ([]domain.KuesionerResult, error) {
	return mockGetAllKuesionersReportFunc(ctx, q)
}

type mockActiveKuesionersHandler struct{}

func (h *mockActiveKuesionersHandler) Handle(ctx context.Context, q ActiveKuesioners.ActiveKuesionerQuery) ([]BankSoaldomain.BankSoalDefault, error) {
	return mockActiveKuesionersFunc(ctx, q)
}

type mockActiveKuesionersSingleHandler struct{}

func (h *mockActiveKuesionersSingleHandler) Handle(ctx context.Context, q ActiveKuesionersSingle.ActiveKuesionerSingleQuery) (*BankSoaldomain.BankSoalDefault, error) {
	return mockActiveKuesionersSingleFunc(ctx, q)
}

type mockSetupUuidKuesionerHandler struct{}

func (h *mockSetupUuidKuesionerHandler) Handle(ctx context.Context, cmd SetupUuidKuesioner.SetupUuidKuesionerCommand) (string, error) {
	return mockSetupUuidKuesionerFunc(ctx, cmd)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[CreateKuesioner.CreateKuesionerCommand, string](&mockCreateKuesionerHandler{})
		_ = mediatr.RegisterRequestHandler[SaveKuesionerJawaban.SaveKuesionerJawabanCommand, string](&mockSaveKuesionerJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[GetKuesionerJawaban.GetKuesionerJawabanByUuidQuery, []domain.KuesionerJawabanDefault](&mockGetKuesionerJawabanHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteKuesioner.DeleteKuesionerCommand, string](&mockDeleteKuesionerHandler{})
		_ = mediatr.RegisterRequestHandler[GetKuesioner.GetKuesionerByUuidQuery, *domain.Kuesioner](&mockGetKuesionerHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllKuesioners.GetAllKuesionersQuery, commondomain.Paged[domain.KuesionerDefault]](&mockGetAllKuesionersHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllKuesionerReport.GetAllKuesionersReportQuery, []domain.KuesionerResult](&mockGetAllKuesionersReportHandler{})
		_ = mediatr.RegisterRequestHandler[ActiveKuesioners.ActiveKuesionerQuery, []BankSoaldomain.BankSoalDefault](&mockActiveKuesionersHandler{})
		_ = mediatr.RegisterRequestHandler[ActiveKuesionersSingle.ActiveKuesionerSingleQuery, *BankSoaldomain.BankSoalDefault](&mockActiveKuesionersSingleHandler{})
		_ = mediatr.RegisterRequestHandler[SetupUuidKuesioner.SetupUuidKuesionerCommand, string](&mockSetupUuidKuesionerHandler{})
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

func TestKuesionerPresentation(t *testing.T) {
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
	presentation.ModuleKuesioner(app)

	token := generateToken("user-123", "simak", "ctxFD49JawrQA")

	t.Run("Setup UUID Kuesioner success", func(t *testing.T) {
		mockSetupUuidKuesionerFunc = func(ctx context.Context, cmd SetupUuidKuesioner.SetupUuidKuesionerCommand) (string, error) {
			return "setup complete", nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kuesioner/setupuuid", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "setup complete", res["message"])
	})

	t.Run("Create Kuesioner success", func(t *testing.T) {
		mockCreateKuesionerFunc = func(ctx context.Context, cmd CreateKuesioner.CreateKuesionerCommand) (string, error) {
			assert.Equal(t, "soal-123", cmd.UuidBankSoal)
			return "new-kuesioner-uuid", nil
		}

		form := url.Values{}
		form.Add("bank_soal", "soal-123")
		form.Add("tanggal", "2026-06-12")

		req := httptest.NewRequest("POST", "/api/v2/kuesioner", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-kuesioner-uuid", res["uuid"])
	})

	t.Run("Save Kuesioner Jawaban success", func(t *testing.T) {
		mockSaveKuesionerJawabanFunc = func(ctx context.Context, cmd SaveKuesionerJawaban.SaveKuesionerJawabanCommand) (string, error) {
			assert.Equal(t, "kues-123", cmd.UuidKuesioner)
			assert.Equal(t, "pert-123", cmd.UuidPertanyaan)
			assert.Equal(t, "A", cmd.Jawaban)
			return "new-jawaban-uuid", nil
		}

		form := url.Values{}
		form.Add("pertanyaan", "pert-123")
		form.Add("jawaban", "A")

		req := httptest.NewRequest("POST", "/api/v2/kuesioner/kues-123/jawaban", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-jawaban-uuid", res["uuid"])
	})

	t.Run("Get Kuesioner Jawaban success", func(t *testing.T) {
		jawaban := "A"
		mockGetKuesionerJawabanFunc = func(ctx context.Context, q GetKuesionerJawaban.GetKuesionerJawabanByUuidQuery) ([]domain.KuesionerJawabanDefault, error) {
			assert.Equal(t, "kues-123", q.Uuid)
			return []domain.KuesionerJawabanDefault{
				{FreeText: &jawaban},
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kuesioner/kues-123/jawaban", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.KuesionerJawabanDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 1)
		assert.Equal(t, "A", *res[0].FreeText)
	})

	t.Run("Delete Kuesioner success", func(t *testing.T) {
		mockDeleteKuesionerFunc = func(ctx context.Context, cmd DeleteKuesioner.DeleteKuesionerCommand) (string, error) {
			assert.Equal(t, "kues-123", cmd.Uuid)
			return "kues-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/kuesioner/kues-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "kues-123", res["uuid"])
	})

	t.Run("Get Kuesioner by UUID success", func(t *testing.T) {
		kuesVal := &domain.Kuesioner{
			IdBankSoal: "soal-123",
		}
		mockGetKuesionerFunc = func(ctx context.Context, q GetKuesioner.GetKuesionerByUuidQuery) (*domain.Kuesioner, error) {
			assert.Equal(t, "kues-123", q.Uuid)
			return kuesVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kuesioner/kues-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.Kuesioner
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "soal-123", res.IdBankSoal)
	})

	t.Run("Get All Kuesioners success", func(t *testing.T) {
		uid := uuid.New()
		items := []domain.KuesionerDefault{
			{UUID: uid},
		}
		mockGetAllKuesionersFunc = func(ctx context.Context, q GetAllKuesioners.GetAllKuesionersQuery) (commondomain.Paged[domain.KuesionerDefault], error) {
			return commondomain.Paged[domain.KuesionerDefault]{
				Data:  items,
				Total: 1,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kuesioners?mode=all", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.KuesionerDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 1)
		assert.Equal(t, uid, res[0].UUID)
	})

	t.Run("Get All Kuesioners Report success", func(t *testing.T) {
		reportVal := []domain.KuesionerResult{
			{Judul: "Report 1"},
		}
		mockGetAllKuesionersReportFunc = func(ctx context.Context, q GetAllKuesionerReport.GetAllKuesionersReportQuery) ([]domain.KuesionerResult, error) {
			assert.Equal(t, "Final Exam", *q.JudulBankSoal)
			return reportVal, nil
		}

		form := url.Values{}
		form.Add("judul", "Final Exam")

		req := httptest.NewRequest("POST", "/api/v2/kuesioners/report", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Get Active Kuesioners success", func(t *testing.T) {
		activeVal := []BankSoaldomain.BankSoalDefault{
			{Judul: "Active Exam"},
		}
		mockActiveKuesionersFunc = func(ctx context.Context, q ActiveKuesioners.ActiveKuesionerQuery) ([]BankSoaldomain.BankSoalDefault, error) {
			assert.Equal(t, "prod-123", *q.Prodi)
			return activeVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kuesioners/active", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []BankSoaldomain.BankSoalDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 1)
		assert.Equal(t, "Active Exam", res[0].Judul)
	})

	t.Run("Get Active Kuesioner Single success", func(t *testing.T) {
		activeSingleVal := &BankSoaldomain.BankSoalDefault{
			Judul: "Single Exam",
		}
		mockActiveKuesionersSingleFunc = func(ctx context.Context, q ActiveKuesionersSingle.ActiveKuesionerSingleQuery) (*BankSoaldomain.BankSoalDefault, error) {
			assert.Equal(t, "prod-123", *q.Prodi)
			assert.Equal(t, "soal-123", q.UUID)
			return activeSingleVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/kuesioners/active/soal-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res BankSoaldomain.BankSoalDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "Single Exam", res.Judul)
	})
}
