package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	Prodidomain "UnpakSiamida/modules/prodi/domain"

	GetAllProdis "UnpakSiamida/modules/prodi/application/GetAllProdis"
)

// =======================================================
// GET /prodis
// =======================================================

// GetAllProdisHandler godoc
// @Summary Get all Prodis
// @Tags Prodi
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Prodidomain.ProdiDefault]
// @Router /prodis [get]
func GetAllProdisHandlerfunc(c *fiber.Ctx) error {
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

	query := GetAllProdis.GetAllProdisQuery{
		Search:        search,
		SearchFilters: filters,
	}

	var adapter commonpresentation.OutputAdapter[Prodidomain.ProdiDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[Prodidomain.ProdiDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[Prodidomain.ProdiDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[Prodidomain.ProdiDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[Prodidomain.ProdiDefault]{}
	}

	result, err := mediatr.Send[
		GetAllProdis.GetAllProdisQuery,
		commondomain.Paged[Prodidomain.ProdiDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func ModuleProdi(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := os.Getenv("WHOAMI_URL")

	app.Get("/api/v2/prodis", commonpresentation.SmartCompress(), GetAllProdisHandlerfunc)
}
