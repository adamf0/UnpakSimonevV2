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
	"UnpakSiamida/modules/fakultas/domain"
	"UnpakSiamida/modules/fakultas/presentation"

	GetAllFakultass "UnpakSiamida/modules/fakultas/application/GetAllFakultass"
)

var (
	mockGetAllFakultassFunc func(ctx context.Context, q GetAllFakultass.GetAllFakultassQuery) (commondomain.Paged[domain.FakultasDefault], error)
	registerOnce            sync.Once
)

type mockGetAllFakultassHandler struct{}

func (h *mockGetAllFakultassHandler) Handle(ctx context.Context, q GetAllFakultass.GetAllFakultassQuery) (commondomain.Paged[domain.FakultasDefault], error) {
	return mockGetAllFakultassFunc(ctx, q)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[GetAllFakultass.GetAllFakultassQuery, commondomain.Paged[domain.FakultasDefault]](&mockGetAllFakultassHandler{})
	})
}

func TestFakultasPresentation(t *testing.T) {
	setupMediatrMocks()

	app := fiber.New()
	presentation.ModuleFakultas(app)

	t.Run("Get All Fakultass success", func(t *testing.T) {
		items := []domain.FakultasDefault{
			{NamaFakultas: "Fakultas Teknik"},
			{NamaFakultas: "Fakultas Hukum"},
		}
		mockGetAllFakultassFunc = func(ctx context.Context, q GetAllFakultass.GetAllFakultassQuery) (commondomain.Paged[domain.FakultasDefault], error) {
			return commondomain.Paged[domain.FakultasDefault]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/fakultass", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res commondomain.Paged[domain.FakultasDefault]
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, "Fakultas Teknik", res.Data[0].NamaFakultas)
	})
}
