package presentation

import (
	"io"
	"net/http/httptest"
	"UnpakSiamida/common/domain"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestPagingAndAllAdapters(t *testing.T) {
	app := fiber.New()

	pagedData := domain.Paged[dummyItem]{
		Data: []dummyItem{
			{ID: 1, Name: "A"},
			{ID: 2, Name: "B"},
		},
		Total:       2,
		CurrentPage: 1,
		TotalPages:  1,
	}

	app.Get("/paging", func(c *fiber.Ctx) error {
		adapter := &PagingAdapter[dummyItem]{}
		return adapter.Send(c, pagedData)
	})

	app.Get("/all", func(c *fiber.Ctx) error {
		adapter := &AllAdapter[dummyItem]{}
		return adapter.Send(c, pagedData)
	})

	// 1) Test PagingAdapter
	reqPaging := httptest.NewRequest("GET", "/paging", nil)
	respPaging, err := app.Test(reqPaging)
	require.NoError(t, err)
	assert.Equal(t, 200, respPaging.StatusCode)

	bodyPaging, _ := io.ReadAll(respPaging.Body)
	assert.Contains(t, string(bodyPaging), `"total":2`)
	assert.Contains(t, string(bodyPaging), `"data":`)

	// 2) Test AllAdapter
	reqAll := httptest.NewRequest("GET", "/all", nil)
	respAll, err := app.Test(reqAll)
	require.NoError(t, err)
	assert.Equal(t, 200, respAll.StatusCode)

	bodyAll, _ := io.ReadAll(respAll.Body)
	assert.NotContains(t, string(bodyAll), `"total":`)
	assert.Contains(t, string(bodyAll), `[{"id":1,"name":"A"},{"id":2,"name":"B"}]`)
}

func TestNDJSONAndSSEAdapters(t *testing.T) {
	app := fiber.New()

	pagedData := domain.Paged[dummyItem]{
		Data: []dummyItem{
			{ID: 1, Name: "A"},
		},
		Total: 1,
	}

	app.Get("/ndjson", func(c *fiber.Ctx) error {
		adapter := &NDJSONAdapter[dummyItem]{}
		return adapter.Send(c, pagedData)
	})

	app.Get("/sse", func(c *fiber.Ctx) error {
		adapter := &SSEAdapter[dummyItem]{}
		return adapter.Send(c, pagedData)
	})

	// 1) Test NDJSONAdapter
	reqNDJSON := httptest.NewRequest("GET", "/ndjson", nil)
	respNDJSON, err := app.Test(reqNDJSON)
	require.NoError(t, err)
	assert.Equal(t, "application/x-ndjson", respNDJSON.Header.Get("Content-Type"))

	bodyNDJSON, _ := io.ReadAll(respNDJSON.Body)
	assert.Contains(t, string(bodyNDJSON), `{"id":1,"name":"A"}`)

	// 2) Test SSEAdapter
	reqSSE := httptest.NewRequest("GET", "/sse", nil)
	respSSE, err := app.Test(reqSSE)
	require.NoError(t, err)
	assert.Equal(t, "text/event-stream", respSSE.Header.Get("Content-Type"))

	bodySSE, _ := io.ReadAll(respSSE.Body)
	assert.Contains(t, string(bodySSE), "total: 1")
	assert.Contains(t, string(bodySSE), "data: start")
	assert.Contains(t, string(bodySSE), `data: {"id":1,"name":"A"}`)
	assert.Contains(t, string(bodySSE), "data: done")
}

func TestWebSocketAdapter_ErrUpgrade(t *testing.T) {
	app := fiber.New()

	app.Get("/ws", func(c *fiber.Ctx) error {
		adapter := &WebSocketAdapter[dummyItem]{}
		return adapter.Send(c, domain.Paged[dummyItem]{})
	})

	req := httptest.NewRequest("GET", "/ws", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 426, resp.StatusCode) // Upgrade Required
}
