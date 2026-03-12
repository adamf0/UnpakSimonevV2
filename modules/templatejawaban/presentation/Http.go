package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	TemplateJawabandomain "UnpakSiamida/modules/templatejawaban/domain"

	CreateTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/CreateTemplateJawaban"
	DeleteTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/DeleteTemplateJawaban"
	GetAllTemplateJawabans "UnpakSiamida/modules/templatejawaban/application/GetAllTemplateJawabans"
	GetTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/GetTemplateJawaban"
	RestoreTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/RestoreTemplateJawaban"
	SetupUuidTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/SetupUuidTemplateJawaban"
	UpdateTemplateJawaban "UnpakSiamida/modules/templatejawaban/application/UpdateTemplateJawaban"
)

// =======================================================
// POST /templatejawaban
// =======================================================

// CreateTemplateJawabanHandler godoc
// @Summary Create new TemplateJawaban
// @Tags TemplateJawaban
//@param template_pertanyaan formData string true "TemplatePertanyaan" format(uuid)
//@param jawaban formData string true "Jawaban"
//@param nilai formData string true "Nilai"
//@param isFreeText formData string true "IsFreeText (0 = false, 1 = true)"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created TemplateJawaban"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatejawaban [post]

func CreateTemplateJawabanHandlerfunc(c *fiber.Ctx) error {

	TemplatePertanyaan := c.FormValue("template_pertanyaan")
	Jawaban := c.FormValue("jawaban")
	Nilai := c.FormValue("nilai")
	IsFreeText := c.FormValue("isFreeText")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := CreateTemplateJawaban.CreateTemplateJawabanCommand{
		UuidTemplatePertanyaan: TemplatePertanyaan,
		Jawaban:                Jawaban,
		Nilai:                  Nilai,
		IsFreeText:             IsFreeText,
		SID:                    SID,
		Resource:               Resource,
	}

	uuid, err := mediatr.Send[CreateTemplateJawaban.CreateTemplateJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /templatejawaban/{uuid}
// =======================================================

// UpdateTemplateJawabanHandler godoc
// @Summary Update existing TemplateJawaban
// @Tags TemplateJawaban
// @Param uuid path string true "TemplateJawaban UUID" format(uuid)
// @param template_pertanyaan formData string true "TemplatePertanyaan" format(uuid)
// @param jawaban formData string true "Jawaban"
// @param nilai formData string true "Nilai"
// @param isFreeText formData string true "IsFreeText (0 = false, 1 = true)"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated TemplateJawaban"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatejawaban/{uuid} [put]
func UpdateTemplateJawabanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	TemplatePertanyaan := c.FormValue("template_pertanyaan")
	Jawaban := c.FormValue("jawaban")
	Nilai := c.FormValue("nilai")
	IsFreeText := c.FormValue("isFreeText")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := UpdateTemplateJawaban.UpdateTemplateJawabanCommand{
		Uuid:                   uuid,
		UuidTemplatePertanyaan: TemplatePertanyaan,
		Jawaban:                Jawaban,
		Nilai:                  Nilai,
		IsFreeText:             IsFreeText,
		SID:                    SID,
		Resource:               Resource,
	}

	updatedID, err := mediatr.Send[UpdateTemplateJawaban.UpdateTemplateJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /templatejawaban/{uuid}
// =======================================================

// DeleteTemplateJawabanHandler godoc
// @Summary Delete a TemplateJawaban
// @Tags TemplateJawaban
// @Param uuid path string true "TemplateJawaban UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplateJawaban"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatejawaban/{uuid} [delete]
func DeleteTemplateJawabanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteTemplateJawaban.DeleteTemplateJawabanCommand{
		Uuid: uuid,
		Mode: "soft_delete",
	}

	deletedID, err := mediatr.Send[DeleteTemplateJawaban.DeleteTemplateJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// DELETE /templatejawaban/{uuid}/force
// =======================================================

// ForceTemplateJawabanHandler godoc
// @Summary Delete a TemplateJawaban
// @Tags TemplateJawaban
// @Param uuid path string true "TemplateJawaban UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplateJawaban"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatejawaban/{uuid} [delete]
func ForceDeleteTemplateJawabanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteTemplateJawaban.DeleteTemplateJawabanCommand{
		Uuid: uuid,
		Mode: "hard_delete",
	}

	deletedID, err := mediatr.Send[DeleteTemplateJawaban.DeleteTemplateJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// PUT /templatejawaban/{uuid}/restore
// =======================================================

// RestoreTemplateJawabanHandler godoc
// @Summary Restore a TemplateJawaban
// @Tags TemplateJawaban
// @Param uuid path string true "TemplateJawaban UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplateJawaban"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatejawaban/{uuid} [delete]
func RestoreTemplateJawabanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := RestoreTemplateJawaban.RestoreTemplateJawabanCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[RestoreTemplateJawaban.RestoreTemplateJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /TemplateJawaban/{uuid}
// =======================================================

// GetTemplateJawabanHandler godoc
// @Summary Get TemplateJawaban by UUID
// @Tags TemplateJawaban
// @Param uuid path string true "TemplateJawaban UUID" format(uuid)
// @Produce json
// @Success 200 {object} TemplateJawabandomain.TemplateJawaban
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /TemplateJawaban/{uuid} [get]
func GetTemplateJawabanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetTemplateJawaban.GetTemplateJawabanByUuidQuery{Uuid: uuid}

	TemplateJawaban, err := mediatr.Send[
		GetTemplateJawaban.GetTemplateJawabanByUuidQuery,
		*TemplateJawabandomain.TemplateJawaban,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if TemplateJawaban == nil {
		return c.Status(404).JSON(fiber.Map{"error": "TemplateJawaban not found"})
	}

	return c.JSON(TemplateJawaban)
}

// =======================================================
// GET /TemplateJawabans
// =======================================================

// GetAllTemplateJawabansHandler godoc
// @Summary Get all TemplateJawabans
// @Tags TemplateJawaban
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[TemplateJawabandomain.TemplateJawabanDefault]
// @Router /TemplateJawabans [get]
func GetAllTemplateJawabansHandlerfunc(c *fiber.Ctx) error {
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
	query := GetAllTemplateJawabans.GetAllTemplateJawabansQuery{
		Search:        search,
		SearchFilters: filters,
		Deleted:       withDeleted,
	}

	var adapter commonpresentation.OutputAdapter[TemplateJawabandomain.TemplateJawabanDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[TemplateJawabandomain.TemplateJawabanDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[TemplateJawabandomain.TemplateJawabanDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[TemplateJawabandomain.TemplateJawabanDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[TemplateJawabandomain.TemplateJawabanDefault]{}
	}

	result, err := mediatr.Send[
		GetAllTemplateJawabans.GetAllTemplateJawabansQuery,
		commondomain.Paged[TemplateJawabandomain.TemplateJawabanDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidTemplateJawabansHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidTemplateJawaban.SetupUuidTemplateJawabanCommand{}

	message, err := mediatr.Send[SetupUuidTemplateJawaban.SetupUuidTemplateJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleTemplateJawaban(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := "http://localhost:3000/whoami"

	app.Get("/templatejawaban/setupuuid", SetupUuidTemplateJawabansHandlerfunc)

	app.Post("/templatejawaban", commonpresentation.JWTMiddleware(), CreateTemplateJawabanHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)
	app.Put("/templatejawaban/:uuid", commonpresentation.JWTMiddleware(), UpdateTemplateJawabanHandlerfunc)

	app.Delete("/templatejawaban/:uuid", commonpresentation.JWTMiddleware(), DeleteTemplateJawabanHandlerfunc)            //soft delete
	app.Delete("/templatejawaban/:uuid/force", commonpresentation.JWTMiddleware(), ForceDeleteTemplateJawabanHandlerfunc) //hanya lpm saja yg hard delete
	app.Put("/templatejawaban/:uuid/restore", commonpresentation.JWTMiddleware(), RestoreTemplateJawabanHandlerfunc)

	app.Get("/templatejawaban/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTemplateJawabanHandlerfunc)
	app.Get("/templatejawabans", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllTemplateJawabansHandlerfunc)
}
