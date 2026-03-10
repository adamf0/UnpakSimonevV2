package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	Kategoridomain "UnpakSiamida/modules/kategori/domain"

	CopyKategori "UnpakSiamida/modules/kategori/application/CopyKategori"
	CreateKategori "UnpakSiamida/modules/kategori/application/CreateKategori"
	DeleteKategori "UnpakSiamida/modules/kategori/application/DeleteKategori"
	GetAllKategoris "UnpakSiamida/modules/kategori/application/GetAllKategoris"
	GetKategori "UnpakSiamida/modules/kategori/application/GetKategori"
	RestoreKategori "UnpakSiamida/modules/kategori/application/RestoreKategori"
	SetupUuidKategori "UnpakSiamida/modules/kategori/application/SetupUuidKategori"
	UpdateKategori "UnpakSiamida/modules/kategori/application/UpdateKategori"
)

// =======================================================
// POST /kategori
// =======================================================

// CreateKategoriHandler godoc
// @Summary Create new Kategori
// @Tags Kategori
//@param nama_kategori formData string true "NamaKategori"
//@param sub_kategori formData string false "SubKategori" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of created Kategori"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kategori [post]

func CreateKategoriHandlerfunc(c *fiber.Ctx) error {

	NamaKategori := c.FormValue("nama_kategori")
	SubKategori := c.FormValue("sub_kategori")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := CreateKategori.CreateKategoriCommand{
		NamaKategori: NamaKategori,
		SubKategori:  helper.StrPtr(SubKategori),
		SID:          SID,
		Resource:     Resource,
	}

	uuid, err := mediatr.Send[CreateKategori.CreateKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /kategori/{uuid}
// =======================================================

// UpdateKategoriHandler godoc
// @Summary Update existing Kategori
// @Tags Kategori
// @Param uuid path string true "Kategori UUID" format(uuid)
// @param nama_kategori formData string true "NamaKategori"
// @param sub_kategori formData string false "SubKategori" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Kategori"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kategori/{uuid} [put]
func UpdateKategoriHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	NamaKategori := c.FormValue("nama_kategori")
	SubKategori := c.FormValue("sub_kategori")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := UpdateKategori.UpdateKategoriCommand{
		Uuid:         uuid,
		NamaKategori: NamaKategori,
		SubKategori:  helper.StrPtr(SubKategori),
		SID:          SID,
		Resource:     Resource,
	}

	updatedID, err := mediatr.Send[UpdateKategori.UpdateKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /kategori/{uuid}
// =======================================================

// DeleteKategoriHandler godoc
// @Summary Delete a Kategori
// @Tags Kategori
// @Param uuid path string true "Kategori UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Kategori"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kategori/{uuid} [delete]
func DeleteKategoriHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteKategori.DeleteKategoriCommand{
		Uuid: uuid,
		Mode: "soft_delete",
	}

	deletedID, err := mediatr.Send[DeleteKategori.DeleteKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// DELETE /kategori/{uuid}/force
// =======================================================

// ForceKategoriHandler godoc
// @Summary Delete a Kategori
// @Tags Kategori
// @Param uuid path string true "Kategori UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Kategori"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kategori/{uuid} [delete]
func ForceDeleteKategoriHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteKategori.DeleteKategoriCommand{
		Uuid: uuid,
		Mode: "hard_delete",
	}

	deletedID, err := mediatr.Send[DeleteKategori.DeleteKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// PUT /kategori/{uuid}/restore
// =======================================================

// RestoreKategoriHandler godoc
// @Summary Restore a Kategori
// @Tags Kategori
// @Param uuid path string true "Kategori UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Kategori"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kategori/{uuid} [delete]
func RestoreKategoriHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := RestoreKategori.RestoreKategoriCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[RestoreKategori.RestoreKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// POST /kategori/{uuid}/copy
// =======================================================

// COpyKategoriHandler godoc
// @Summary Copy a Kategori
// @Tags Kategori
// @Param uuid path string true "Kategori UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Kategori"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kategori/{uuid} [delete]
func CopyKategoriHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := CopyKategori.CopyKategoriCommand{
		Uuid:     uuid,
		SID:      SID,
		Resource: Resource,
	}

	deletedID, err := mediatr.Send[CopyKategori.CopyKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /Kategori/{uuid}
// =======================================================

// GetKategoriHandler godoc
// @Summary Get Kategori by UUID
// @Tags Kategori
// @Param uuid path string true "Kategori UUID" format(uuid)
// @Produce json
// @Success 200 {object} Kategoridomain.Kategori
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /Kategori/{uuid} [get]
func GetKategoriHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetKategori.GetKategoriByUuidQuery{Uuid: uuid}

	Kategori, err := mediatr.Send[
		GetKategori.GetKategoriByUuidQuery,
		*Kategoridomain.Kategori,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if Kategori == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kategori not found"})
	}

	return c.JSON(Kategori)
}

// =======================================================
// GET /Kategoris
// =======================================================

// GetAllKategorisHandler godoc
// @Summary Get all Kategoris
// @Tags Kategori
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Kategoridomain.KategoriDefault]
// @Router /Kategoris [get]
func GetAllKategorisHandlerfunc(c *fiber.Ctx) error {
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
	query := GetAllKategoris.GetAllKategorisQuery{
		Search:        search,
		SearchFilters: filters,
		Deleted:       withDeleted,
	}

	var adapter commonpresentation.OutputAdapter[Kategoridomain.KategoriDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[Kategoridomain.KategoriDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[Kategoridomain.KategoriDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[Kategoridomain.KategoriDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[Kategoridomain.KategoriDefault]{}
	}

	result, err := mediatr.Send[
		GetAllKategoris.GetAllKategorisQuery,
		commondomain.Paged[Kategoridomain.KategoriDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidKategorisHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidKategori.SetupUuidKategoriCommand{}

	message, err := mediatr.Send[SetupUuidKategori.SetupUuidKategoriCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleKategori(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := "http://localhost:3000/whoami"

	app.Get("/kategori/setupuuid", SetupUuidKategorisHandlerfunc)

	app.Post("/kategori", commonpresentation.JWTMiddleware(), CreateKategoriHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)
	app.Put("/kategori/:uuid", commonpresentation.JWTMiddleware(), UpdateKategoriHandlerfunc)

	app.Delete("/kategori/:uuid", commonpresentation.JWTMiddleware(), DeleteKategoriHandlerfunc)            //soft delete
	app.Delete("/kategori/:uuid/force", commonpresentation.JWTMiddleware(), ForceDeleteKategoriHandlerfunc) //hanya lpm saja yg hard delete
	app.Put("/kategori/:uuid/restore", commonpresentation.JWTMiddleware(), RestoreKategoriHandlerfunc)
	app.Post("/kategori/:uuid/copy", commonpresentation.JWTMiddleware(), CopyKategoriHandlerfunc)

	app.Get("/kategori/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetKategoriHandlerfunc)
	app.Get("/kategoris", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllKategorisHandlerfunc)
}
