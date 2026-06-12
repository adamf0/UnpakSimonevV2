package presentation

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRawHttpRequest(t *testing.T) {
	app := fiber.New()
	app.Post("/test", func(c *fiber.Ctx) error {
		duration := 150 * time.Millisecond
		raw := BuildRawHttpRequest(c, duration)

		assert.Contains(t, raw, "POST /test?param=val HTTP/1.1")
		assert.Contains(t, raw, "X-Custom-Header: HeaderValue")
		assert.Contains(t, raw, "my-body-content")
		assert.Contains(t, raw, "Took: 150ms")
		return c.SendString("OK")
	})

	req := httptest.NewRequest("POST", "/test?param=val", strings.NewReader("my-body-content"))
	req.Header.Set("X-Custom-Header", "HeaderValue")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLoggerMiddleware(t *testing.T) {
	// Remove existing request.log if any
	_ = os.Remove("request.log")

	app := fiber.New()
	app.Use(LoggerMiddleware)
	app.Get("/test-logger", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test-logger", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify request.log was created and contains the request
	_, errStat := os.Stat("request.log")
	assert.NoError(t, errStat)

	content, errRead := os.ReadFile("request.log")
	require.NoError(t, errRead)
	assert.Contains(t, string(content), "GET /test-logger HTTP/1.1")

	// Clean up
	_ = os.Remove("request.log")
}
