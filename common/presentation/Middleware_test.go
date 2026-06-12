package presentation

import (
	"UnpakSiamida/common/helper"
	"io"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTMiddleware_Success(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddleware())
	app.Get("/secure", func(c *fiber.Ctx) error {
		sid := string(c.Request().PostArgs().Peek("sid"))
		resource := string(c.Request().PostArgs().Peek("resource"))
		codectx := string(c.Request().PostArgs().Peek("codectx"))
		return c.JSON(fiber.Map{
			"sid":      sid,
			"resource": resource,
			"codectx":  codectx,
		})
	})

	// Generate valid token
	codectxVal := "dosen"
	accessToken, _, err := helper.GenerateToken("sid-123", "simak", &codectxVal)
	require.NoError(t, err)

	// Test passing via Header
	req := httptest.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"sid":"sid-123"`)
	assert.Contains(t, string(body), `"resource":"simak"`)
	assert.Contains(t, string(body), `"codectx":"dosen"`)

	// Test passing via Query param
	reqQuery := httptest.NewRequest("GET", "/secure?ctxtoken=Bearer%20"+accessToken, nil)
	respQuery, err := app.Test(reqQuery)
	require.NoError(t, err)
	assert.Equal(t, 200, respQuery.StatusCode)
}

func TestJWTMiddleware_FailMissingOrInvalid(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddleware())
	app.Get("/secure", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Missing header
	req := httptest.NewRequest("GET", "/secure", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Invalid format
	reqInvalid := httptest.NewRequest("GET", "/secure", nil)
	reqInvalid.Header.Set("Authorization", "InvalidFormat token")
	respInvalid, err := app.Test(reqInvalid)
	require.NoError(t, err)
	assert.Equal(t, 400, respInvalid.StatusCode)
}

func TestSmartCompress(t *testing.T) {
	app := fiber.New()
	app.Use(SmartCompress())

	app.Get("/text", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.SendString(strings.Repeat("Some text message that is long enough to be compressed. ", 100))
	})

	app.Get("/stream", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		return c.SendString("Some stream data")
	})

	// Test text/plain compress header
	reqText := httptest.NewRequest("GET", "/text", nil)
	reqText.Header.Set("Accept-Encoding", "gzip")
	respText, err := app.Test(reqText)
	require.NoError(t, err)
	assert.Equal(t, "gzip", respText.Header.Get("Content-Encoding"))

	// Test text/event-stream not compressed
	reqStream := httptest.NewRequest("GET", "/stream", nil)
	reqStream.Header.Set("Accept-Encoding", "gzip")
	respStream, err := app.Test(reqStream)
	require.NoError(t, err)
	assert.Empty(t, respStream.Header.Get("Content-Encoding"))
}
