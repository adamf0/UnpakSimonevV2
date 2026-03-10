package presentation

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	TemplatePertanyaandomain "UnpakSiamida/modules/templatepertanyaan/domain"

	CopyTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/CopyTemplatePertanyaan"
	CreateTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/CreateTemplatePertanyaan"
	DeleteTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/DeleteTemplatePertanyaan"
	GetAllTemplatePertanyaans "UnpakSiamida/modules/templatepertanyaan/application/GetAllTemplatePertanyaans"
	GetTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaan"
	RestoreTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/RestoreTemplatePertanyaan"
	SetupUuidTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/SetupUuidTemplatePertanyaan"
	UpdateTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/application/UpdateTemplatePertanyaan"
)

// =======================================================
// POST /templatepertanyaan
// =======================================================

// CreateTemplatePertanyaanHandler godoc
// @Summary Create new TemplatePertanyaan
// @Tags TemplatePertanyaan
//@param bank_soal formData string true "BankSoal" format(uuid)
//@param pertanyaan formData string true "Pertanyaan"
//@param jenis_pilihan formData string true "JenisPilihan"
//@param bobot formData string true "Bobot"
//@param kategori formData string true "Kategori" format(uuid)
//@param required formData int true "Required (0 = false, 1 = true)"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created TemplatePertanyaan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatepertanyaan [post]

func CreateTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {

	BankSoal := c.FormValue("bank_soal")
	Pertanyaan := c.FormValue("pertanyaan")
	JenisPilihan := c.FormValue("jenis_pilihan")
	Bobot := c.FormValue("bobot")
	Kategori := c.FormValue("kategori")
	Required := c.FormValue("required")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	required := 0
	if v, err := strconv.ParseBool(Required); err == nil && v {
		required = 1
	}

	cmd := CreateTemplatePertanyaan.CreateTemplatePertanyaanCommand{
		UuidBankSoal: BankSoal,
		Pertanyaan:   Pertanyaan,
		JenisPilihan: JenisPilihan,
		Bobot:        Bobot,
		UuidKategori: Kategori,
		Required:     required,
		SID:          SID,
		Resource:     Resource,
	}

	uuid, err := mediatr.Send[CreateTemplatePertanyaan.CreateTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /templatepertanyaan/{uuid}
// =======================================================

// UpdateTemplatePertanyaanHandler godoc
// @Summary Update existing TemplatePertanyaan
// @Tags TemplatePertanyaan
// @Param uuid path string true "TemplatePertanyaan UUID" format(uuid)
// @param bank_soal formData string true "BankSoal" format(uuid)
// @param pertanyaan formData string true "Pertanyaan"
// @param jenis_pilihan formData string true "JenisPilihan"
// @param bobot formData string true "Bobot"
// @param kategori formData string true "Kategori" format(uuid)
// @param required formData int true "Required (0 = false, 1 = true)"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated TemplatePertanyaan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatepertanyaan/{uuid} [put]
func UpdateTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	BankSoal := c.FormValue("bank_soal")
	Pertanyaan := c.FormValue("pertanyaan")
	JenisPilihan := c.FormValue("jenis_pilihan")
	Bobot := c.FormValue("bobot")
	Kategori := c.FormValue("kategori")
	Required := c.FormValue("required")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	required := 0
	if v, err := strconv.ParseBool(Required); err == nil && v {
		required = 1
	}

	cmd := UpdateTemplatePertanyaan.UpdateTemplatePertanyaanCommand{
		Uuid:         uuid,
		UuidBankSoal: BankSoal,
		Pertanyaan:   Pertanyaan,
		JenisPilihan: JenisPilihan,
		Bobot:        Bobot,
		UuidKategori: Kategori,
		Required:     required,
		SID:          SID,
		Resource:     Resource,
	}

	updatedID, err := mediatr.Send[UpdateTemplatePertanyaan.UpdateTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /templatepertanyaan/{uuid}
// =======================================================

// DeleteTemplatePertanyaanHandler godoc
// @Summary Delete a TemplatePertanyaan
// @Tags TemplatePertanyaan
// @Param uuid path string true "TemplatePertanyaan UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplatePertanyaan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatepertanyaan/{uuid} [delete]
func DeleteTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand{
		Uuid: uuid,
		Mode: "soft_delete",
	}

	deletedID, err := mediatr.Send[DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// DELETE /templatepertanyaan/{uuid}/force
// =======================================================

// ForceTemplatePertanyaanHandler godoc
// @Summary Delete a TemplatePertanyaan
// @Tags TemplatePertanyaan
// @Param uuid path string true "TemplatePertanyaan UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplatePertanyaan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatepertanyaan/{uuid} [delete]
func ForceDeleteTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand{
		Uuid: uuid,
		Mode: "hard_delete",
	}

	deletedID, err := mediatr.Send[DeleteTemplatePertanyaan.DeleteTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// PUT /templatepertanyaan/{uuid}/restore
// =======================================================

// RestoreTemplatePertanyaanHandler godoc
// @Summary Restore a TemplatePertanyaan
// @Tags TemplatePertanyaan
// @Param uuid path string true "TemplatePertanyaan UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplatePertanyaan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatepertanyaan/{uuid} [delete]
func RestoreTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := RestoreTemplatePertanyaan.RestoreTemplatePertanyaanCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[RestoreTemplatePertanyaan.RestoreTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// POST /templatepertanyaan/{uuid}/copy
// =======================================================

// COpyTemplatePertanyaanHandler godoc
// @Summary Copy a TemplatePertanyaan
// @Tags TemplatePertanyaan
// @Param uuid path string true "TemplatePertanyaan UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted TemplatePertanyaan"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /templatepertanyaan/{uuid} [delete]
func CopyTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := CopyTemplatePertanyaan.CopyTemplatePertanyaanCommand{
		Uuid:     uuid,
		SID:      SID,
		Resource: Resource,
	}

	deletedID, err := mediatr.Send[CopyTemplatePertanyaan.CopyTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /TemplatePertanyaan/{uuid}
// =======================================================

// GetTemplatePertanyaanHandler godoc
// @Summary Get TemplatePertanyaan by UUID
// @Tags TemplatePertanyaan
// @Param uuid path string true "TemplatePertanyaan UUID" format(uuid)
// @Produce json
// @Success 200 {object} TemplatePertanyaandomain.TemplatePertanyaan
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /TemplatePertanyaan/{uuid} [get]
func GetTemplatePertanyaanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetTemplatePertanyaan.GetTemplatePertanyaanByUuidQuery{Uuid: uuid}

	TemplatePertanyaan, err := mediatr.Send[
		GetTemplatePertanyaan.GetTemplatePertanyaanByUuidQuery,
		*TemplatePertanyaandomain.TemplatePertanyaan,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if TemplatePertanyaan == nil {
		return c.Status(404).JSON(fiber.Map{"error": "TemplatePertanyaan not found"})
	}

	return c.JSON(TemplatePertanyaan)
}

// =======================================================
// GET /TemplatePertanyaans
// =======================================================

// GetAllTemplatePertanyaansHandler godoc
// @Summary Get all TemplatePertanyaans
// @Tags TemplatePertanyaan
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[TemplatePertanyaandomain.TemplatePertanyaanDefault]
// @Router /TemplatePertanyaans [get]
func GetAllTemplatePertanyaansHandlerfunc(c *fiber.Ctx) error {
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
	query := GetAllTemplatePertanyaans.GetAllTemplatePertanyaansQuery{
		Search:        search,
		SearchFilters: filters,
		Deleted:       withDeleted,
	}

	var adapter commonpresentation.OutputAdapter[TemplatePertanyaandomain.TemplatePertanyaanDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[TemplatePertanyaandomain.TemplatePertanyaanDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[TemplatePertanyaandomain.TemplatePertanyaanDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[TemplatePertanyaandomain.TemplatePertanyaanDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[TemplatePertanyaandomain.TemplatePertanyaanDefault]{}
	}

	result, err := mediatr.Send[
		GetAllTemplatePertanyaans.GetAllTemplatePertanyaansQuery,
		commondomain.Paged[TemplatePertanyaandomain.TemplatePertanyaanDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidTemplatePertanyaansHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidTemplatePertanyaan.SetupUuidTemplatePertanyaanCommand{}

	message, err := mediatr.Send[SetupUuidTemplatePertanyaan.SetupUuidTemplatePertanyaanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleTemplatePertanyaan(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := "http://localhost:3000/whoami"

	app.Get("/templatepertanyaan/setupuuid", SetupUuidTemplatePertanyaansHandlerfunc)

	app.Post("/templatepertanyaan", commonpresentation.JWTMiddleware(), CreateTemplatePertanyaanHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)
	app.Put("/templatepertanyaan/:uuid", commonpresentation.JWTMiddleware(), UpdateTemplatePertanyaanHandlerfunc)

	app.Delete("/templatepertanyaan/:uuid", commonpresentation.JWTMiddleware(), DeleteTemplatePertanyaanHandlerfunc)            //soft delete
	app.Delete("/templatepertanyaan/:uuid/force", commonpresentation.JWTMiddleware(), ForceDeleteTemplatePertanyaanHandlerfunc) //hanya lpm saja yg hard delete
	app.Put("/templatepertanyaan/:uuid/restore", commonpresentation.JWTMiddleware(), RestoreTemplatePertanyaanHandlerfunc)
	app.Post("/templatepertanyaan/:uuid/copy", commonpresentation.JWTMiddleware(), CopyTemplatePertanyaanHandlerfunc)

	app.Get("/templatepertanyaan/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetTemplatePertanyaanHandlerfunc)
	app.Get("/templatepertanyaans", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllTemplatePertanyaansHandlerfunc)
}
