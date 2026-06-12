package presentation

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonUUID(t *testing.T) {
	app := fiber.New()
	app.Get("/uuid", func(c *fiber.Ctx) error {
		return JsonUUID(c, "550e8400-e29b-41d4-a716-446655440000")
	})

	req := httptest.NewRequest("GET", "/uuid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"uuid":"550e8400-e29b-41d4-a716-446655440000"`)
}
