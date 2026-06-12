package infrastructure_test

import (
	"errors"
	"net/http/httptest"
	"UnpakSiamida/common/domain"
	"UnpakSiamida/common/infrastructure"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseError_Format(t *testing.T) {
	err := infrastructure.NewResponseError("FAIL.CODE", "Failure occurred")
	assert.Equal(t, "Failure occurred", err.Error())

	errMap := infrastructure.NewResponseError("FAIL.MAP", map[string]string{"err": "details"})
	assert.Contains(t, errMap.Error(), "FAIL.MAP:")
}

func TestResponseError_HandleError(t *testing.T) {
	app := fiber.New()

	app.Get("/error", func(c *fiber.Ctx) error {
		errType := c.Query("type")
		switch errType {
		case "response":
			return infrastructure.HandleError(c, infrastructure.NewResponseError("Test.NotFound", "resource missing"))
		case "domain":
			derr := domain.NotFoundError("Domain.NotFound", "domain resource missing")
			return infrastructure.HandleError(c, derr)
		case "internal":
			return infrastructure.HandleError(c, errors.New("raw db error"))
		default:
			return infrastructure.HandleError(c, nil)
		}
	})

	// 1) Test Nil error
	reqNil := httptest.NewRequest("GET", "/error", nil)
	respNil, err := app.Test(reqNil)
	require.NoError(t, err)
	assert.Equal(t, 200, respNil.StatusCode)

	// 2) Test ResponseError (NotFound status 404)
	reqResp := httptest.NewRequest("GET", "/error?type=response", nil)
	respResp, err := app.Test(reqResp)
	require.NoError(t, err)
	assert.Equal(t, 404, respResp.StatusCode)

	// 3) Test Domain error (NotFound status 404)
	reqDomain := httptest.NewRequest("GET", "/error?type=domain", nil)
	respDomain, err := app.Test(reqDomain)
	require.NoError(t, err)
	assert.Equal(t, 404, respDomain.StatusCode)

	// 4) Test Generic internal error (status 500)
	reqInternal := httptest.NewRequest("GET", "/error?type=internal", nil)
	respInternal, err := app.Test(reqInternal)
	require.NoError(t, err)
	assert.Equal(t, 500, respInternal.StatusCode)
}
