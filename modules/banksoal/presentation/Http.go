package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	BankSoaldomain "UnpakSiamida/modules/banksoal/domain"

	CopyBankSoal "UnpakSiamida/modules/banksoal/application/CopyBankSoal"
	CreateBankSoal "UnpakSiamida/modules/banksoal/application/CreateBankSoal"
	DeleteBankSoal "UnpakSiamida/modules/banksoal/application/DeleteBankSoal"
	GetAllBankSoals "UnpakSiamida/modules/banksoal/application/GetAllBankSoals"
	GetBankSoal "UnpakSiamida/modules/banksoal/application/GetBankSoal"
	RestoreBankSoal "UnpakSiamida/modules/banksoal/application/RestoreBankSoal"
	SetupUuidBankSoal "UnpakSiamida/modules/banksoal/application/SetupUuidBankSoal"
	UpdateBankSoal "UnpakSiamida/modules/banksoal/application/UpdateBankSoal"
)

// =======================================================
// POST /banksoal
// =======================================================

// CreateBankSoalHandler godoc
// @Summary Create new BankSoal
// @Tags BankSoal
//@param judul formData string true "Judul"
//@param content formData string true "Content"
//@param deskripsi formData string true "Deskripsi"
//@param semester formData string true "Semester"
//@param tanggal_mulai formData string true "TanggalMulai"
//@param tanggal_akhir formData string true "TanggalAkhir"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created BankSoal"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /banksoal [post]

func CreateBankSoalHandlerfunc(c *fiber.Ctx) error {

	Judul := c.FormValue("judul")
	Content := c.FormValue("content")
	Deskripsi := c.FormValue("deskripsi")
	Semester := c.FormValue("semester")
	TanggalMulai := c.FormValue("tanggal_mulai")
	TanggalAkhir := c.FormValue("tanggal_akhir")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := CreateBankSoal.CreateBankSoalCommand{
		Judul:        Judul,
		Content:      Content,
		Deskripsi:    Deskripsi,
		Semester:     Semester,
		TanggalMulai: TanggalMulai,
		TanggalAkhir: TanggalAkhir,
		SID:          SID,
		Resource:     Resource,
	}

	uuid, err := mediatr.Send[CreateBankSoal.CreateBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /banksoal/{uuid}
// =======================================================

// UpdateBankSoalHandler godoc
// @Summary Update existing BankSoal
// @Tags BankSoal
// @Param uuid path string true "BankSoal UUID" format(uuid)
// @param judul formData string true "Judul"
// @param content formData string true "Content"
// @param deskripsi formData string true "Deskripsi"
// @param semester formData string true "Semester"
// @param tanggal_mulai formData string true "TanggalMulai"
// @param tanggal_akhir formData string true "TanggalAkhir"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated BankSoal"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /banksoal/{uuid} [put]
func UpdateBankSoalHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	Judul := c.FormValue("judul")
	Content := c.FormValue("content")
	Deskripsi := c.FormValue("deskripsi")
	Semester := c.FormValue("semester")
	TanggalMulai := c.FormValue("tanggal_mulai")
	TanggalAkhir := c.FormValue("tanggal_akhir")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := UpdateBankSoal.UpdateBankSoalCommand{
		Uuid:         uuid,
		Judul:        Judul,
		Content:      Content,
		Deskripsi:    Deskripsi,
		Semester:     Semester,
		TanggalMulai: TanggalMulai,
		TanggalAkhir: TanggalAkhir,
		SID:          SID,
		Resource:     Resource,
	}

	updatedID, err := mediatr.Send[UpdateBankSoal.UpdateBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /banksoal/{uuid}
// =======================================================

// DeleteBankSoalHandler godoc
// @Summary Delete a BankSoal
// @Tags BankSoal
// @Param uuid path string true "BankSoal UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted BankSoal"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /banksoal/{uuid} [delete]
func DeleteBankSoalHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteBankSoal.DeleteBankSoalCommand{
		Uuid: uuid,
		Mode: "soft_delete",
	}

	deletedID, err := mediatr.Send[DeleteBankSoal.DeleteBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// DELETE /banksoal/{uuid}/force
// =======================================================

// ForceBankSoalHandler godoc
// @Summary Delete a BankSoal
// @Tags BankSoal
// @Param uuid path string true "BankSoal UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted BankSoal"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /banksoal/{uuid} [delete]
func ForceDeleteBankSoalHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteBankSoal.DeleteBankSoalCommand{
		Uuid: uuid,
		Mode: "hard_delete",
	}

	deletedID, err := mediatr.Send[DeleteBankSoal.DeleteBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// PUT /banksoal/{uuid}/restore
// =======================================================

// RestoreBankSoalHandler godoc
// @Summary Restore a BankSoal
// @Tags BankSoal
// @Param uuid path string true "BankSoal UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted BankSoal"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /banksoal/{uuid} [delete]
func RestoreBankSoalHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := RestoreBankSoal.RestoreBankSoalCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[RestoreBankSoal.RestoreBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// POST /banksoal/{uuid}/copy
// =======================================================

// CopyBankSoalHandler godoc
// @Summary Copy a BankSoal
// @Tags BankSoal
// @Param uuid path string true "BankSoal UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted BankSoal"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /banksoal/{uuid} [delete]
func CopyBankSoalHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")

	cmd := CopyBankSoal.CopyBankSoalCommand{
		Uuid:     uuid,
		SID:      SID,
		Resource: Resource,
	}

	deletedID, err := mediatr.Send[CopyBankSoal.CopyBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /BankSoal/{uuid}
// =======================================================

// GetBankSoalHandler godoc
// @Summary Get BankSoal by UUID
// @Tags BankSoal
// @Param uuid path string true "BankSoal UUID" format(uuid)
// @Produce json
// @Success 200 {object} BankSoaldomain.BankSoal
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /BankSoal/{uuid} [get]
func GetBankSoalHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetBankSoal.GetBankSoalByUuidQuery{Uuid: uuid}

	BankSoal, err := mediatr.Send[
		GetBankSoal.GetBankSoalByUuidQuery,
		*BankSoaldomain.BankSoal,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if BankSoal == nil {
		return c.Status(404).JSON(fiber.Map{"error": "BankSoal not found"})
	}

	return c.JSON(BankSoal)
}

// =======================================================
// GET /BankSoals
// =======================================================

// GetAllBankSoalsHandler godoc
// @Summary Get all BankSoals
// @Tags BankSoal
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[BankSoaldomain.BankSoalDefault]
// @Router /BankSoals [get]
func GetAllBankSoalsHandlerfunc(c *fiber.Ctx) error {
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
	query := GetAllBankSoals.GetAllBankSoalsQuery{
		Search:        search,
		SearchFilters: filters,
		Deleted:       withDeleted,
	}

	var adapter commonpresentation.OutputAdapter[BankSoaldomain.BankSoalDefault]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[BankSoaldomain.BankSoalDefault]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[BankSoaldomain.BankSoalDefault]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[BankSoaldomain.BankSoalDefault]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[BankSoaldomain.BankSoalDefault]{}
	}

	result, err := mediatr.Send[
		GetAllBankSoals.GetAllBankSoalsQuery,
		commondomain.Paged[BankSoaldomain.BankSoalDefault],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func SetupUuidBankSoalsHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidBankSoal.SetupUuidBankSoalCommand{}

	message, err := mediatr.Send[SetupUuidBankSoal.SetupUuidBankSoalCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleBankSoal(app *fiber.App) {
	// admin := []string{"admin"}
	// whoamiURL := "http://localhost:3000/whoami"

	app.Get("/banksoal/setupuuid", SetupUuidBankSoalsHandlerfunc)

	app.Post("/banksoal", commonpresentation.JWTMiddleware(), CreateBankSoalHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)
	app.Put("/banksoal/:uuid", commonpresentation.JWTMiddleware(), UpdateBankSoalHandlerfunc)

	app.Delete("/banksoal/:uuid", commonpresentation.JWTMiddleware(), DeleteBankSoalHandlerfunc)            //soft delete
	app.Delete("/banksoal/:uuid/force", commonpresentation.JWTMiddleware(), ForceDeleteBankSoalHandlerfunc) //hanya lpm saja yg hard delete
	app.Put("/banksoal/:uuid/restore", commonpresentation.JWTMiddleware(), RestoreBankSoalHandlerfunc)
	app.Post("/banksoal/:uuid/copy", commonpresentation.JWTMiddleware(), CopyBankSoalHandlerfunc)

	app.Get("/banksoal/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetBankSoalHandlerfunc)
	app.Get("/banksoals", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllBankSoalsHandlerfunc)
}
