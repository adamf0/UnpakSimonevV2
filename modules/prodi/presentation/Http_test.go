package presentation_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/prodi/domain"
	"UnpakSiamida/modules/prodi/presentation"

	GetAllProdis "UnpakSiamida/modules/prodi/application/GetAllProdis"
)

var (
	mockGetAllProdisFunc func(ctx context.Context, q GetAllProdis.GetAllProdisQuery) (commondomain.Paged[domain.ProdiDefault], error)
	registerOnce         sync.Once
)

type mockGetAllProdisHandler struct{}

func (h *mockGetAllProdisHandler) Handle(ctx context.Context, q GetAllProdis.GetAllProdisQuery) (commondomain.Paged[domain.ProdiDefault], error) {
	return mockGetAllProdisFunc(ctx, q)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[GetAllProdis.GetAllProdisQuery, commondomain.Paged[domain.ProdiDefault]](&mockGetAllProdisHandler{})
	})
}

func TestProdiPresentation(t *testing.T) {
	setupMediatrMocks()

	app := fiber.New()
	presentation.ModuleProdi(app)

	t.Run("Get All Prodis success", func(t *testing.T) {
		items := []domain.ProdiDefault{
			{NamaProdi: "Ilmu Komputer"},
			{NamaProdi: "Manajemen"},
		}
		mockGetAllProdisFunc = func(ctx context.Context, q GetAllProdis.GetAllProdisQuery) (commondomain.Paged[domain.ProdiDefault], error) {
			return commondomain.Paged[domain.ProdiDefault]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/prodis", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res commondomain.Paged[domain.ProdiDefault]
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, "Ilmu Komputer", res.Data[0].NamaProdi)
	})
}
