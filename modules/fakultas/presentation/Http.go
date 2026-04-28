package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	Fakultasdomain "UnpakSiamida/modules/fakultas/domain"

	GetAllFakultass "UnpakSiamida/modules/fakultas/application/GetAllFakultass"
)

// =======================================================
// GET /fakultass
// =======================================================

// GetAllFakultassHandler godoc
// @Summary Get all Fakultass
// @Tags Fakultas
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Fakultasdomain.FakultasDefault]
// @Router /fakultass [get]
func GetAllFakultassHandlerfunc(c *fiber.Ctx) error {
	mode := c.Query("mode", "paging")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	filtersRaw := c.Query("filters", "")
	var filters []commondomain.SearchFilter

	if filtersRaw != "" {
		parts := strings.Split(filtersRaw, ";")
		for _, p := range parts {
			tokens := strings.SplitN(p, ":", 3)
			if len(tokens) != 3 {
				continue
			}

			field := strings.TrimSpace(tokens[0])
			op := strings.TrimSpace(tokens[1])
			rawValue := strings.TrimSpace(tokens[2])

			var ptr *string
			if rawValue != "" && rawValue != "null" {
				ptr = &rawValue
			}

			filters = append(filters, commondomain.SearchFilter{
				Field:    field,
				Operator: op,
				Value:    ptr,
			})
		}
	}

	query := GetAllFakultass.GetAllFakultassQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[Fakultasdomain.FakultasDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[Fakultasdomain.FakultasDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[Fakultasdomain.FakultasDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[Fakultasdomain.FakultasDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[Fakultasdomain.FakultasDefault]{}
	}

	result, err := mediatr.Send[
		GetAllFakultass.GetAllFakultassQuery,
		commondomain.Paged[Fakultasdomain.FakultasDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func ModuleFakultas(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := os.Getenv("WHOAMI_URL")

	app.Get("/api/v2/fakultass", commonpresentation.SmartCompress(), GetAllFakultassHandlerfunc)
}
