package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	Kuesionerdomain "UnpakSiamida/modules/kuesioner/domain"

	CreateKuesioner "UnpakSiamida/modules/kuesioner/application/CreateKuesioner"
	DeleteKuesioner "UnpakSiamida/modules/kuesioner/application/DeleteKuesioner"
	GetAllKuesioners "UnpakSiamida/modules/kuesioner/application/GetAllKuesioners"
	GetKuesioner "UnpakSiamida/modules/kuesioner/application/GetKuesioner"
	SetupUuidKuesioner "UnpakSiamida/modules/kuesioner/application/SetupUuidKuesioner"
)

// =======================================================
// POST /kuesioner
// =======================================================

// CreateKuesionerHandler godoc
// @Summary Create new Kuesioner
// @Tags Kuesioner
//@param bank_soal formData string true "BankSoal" format(uuid)
//@param tanggal formData string true "Tanggal"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created Kuesioner"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kuesioner [post]

func CreateKuesionerHandlerfunc(c *fiber.Ctx) error {

	UuidBankSoal := c.FormValue("bank_soal")
	Tanggal := c.FormValue("tanggal")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")
	CodeCtx := c.FormValue("codectx")

	cmd := CreateKuesioner.CreateKuesionerCommand{
		UuidBankSoal: UuidBankSoal,
		Tanggal:      Tanggal,
		SID:          SID,
		Resource:     Resource,
		CodeCtx:      CodeCtx,
	}

	uuid, err := mediatr.Send[CreateKuesioner.CreateKuesionerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// DELETE /kuesioner/{uuid}
// =======================================================

// DeleteKuesionerHandler godoc
// @Summary Delete a Kuesioner
// @Tags Kuesioner
// @Param uuid path string true "Kuesioner UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Kuesioner"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kuesioner/{uuid} [delete]
func DeleteKuesionerHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteKuesioner.DeleteKuesionerCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[DeleteKuesioner.DeleteKuesionerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /Kuesioner/{uuid}
// =======================================================

// GetKuesionerHandler godoc
// @Summary Get Kuesioner by UUID
// @Tags Kuesioner
// @Param uuid path string true "Kuesioner UUID" format(uuid)
// @Produce json
// @Success 200 {object} Kuesionerdomain.Kuesioner
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /Kuesioner/{uuid} [get]
func GetKuesionerHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetKuesioner.GetKuesionerByUuidQuery{Uuid: uuid}

	Kuesioner, err := mediatr.Send[
		GetKuesioner.GetKuesionerByUuidQuery,
		*Kuesionerdomain.Kuesioner,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if Kuesioner == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kuesioner not found"})
	}

	return c.JSON(Kuesioner)
}

// =======================================================
// GET /Kuesioners
// =======================================================

// GetAllKuesionersHandler godoc
// @Summary Get all Kuesioners
// @Tags Kuesioner
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Kuesionerdomain.KuesionerDefault]
// @Router /Kuesioners [get]
func GetAllKuesionersHandlerfunc(c *fiber.Ctx) error {
	flag := c.Query("flag", "none") //with deleted
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

	withDeleted := false
	if flag == "deleted" {
		withDeleted = true
	}
	query := GetAllKuesioners.GetAllKuesionersQuery{
		Search:        search,
		SearchFilters: filters,
		Deleted:       withDeleted,
	}

	var adapter commonpresentation.OutputAdapter[Kuesionerdomain.KuesionerDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[Kuesionerdomain.KuesionerDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[Kuesionerdomain.KuesionerDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[Kuesionerdomain.KuesionerDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[Kuesionerdomain.KuesionerDefault]{}
	}

	result, err := mediatr.Send[
		GetAllKuesioners.GetAllKuesionersQuery,
		commondomain.Paged[Kuesionerdomain.KuesionerDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidKuesionersHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidKuesioner.SetupUuidKuesionerCommand{}

	message, err := mediatr.Send[SetupUuidKuesioner.SetupUuidKuesionerCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleKuesioner(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := "http://localhost:3000/whoami"

	app.Get("/kuesioner/setupuuid", SetupUuidKuesionersHandlerfunc)

	app.Post("/kuesioner", commonpresentation.JWTMiddleware(), CreateKuesionerHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)

	app.Delete("/kuesioner/:uuid", commonpresentation.JWTMiddleware(), DeleteKuesionerHandlerfunc)

	app.Get("/kuesioner/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetKuesionerHandlerfunc)
	app.Get("/kuesioners", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllKuesionersHandlerfunc)
}
