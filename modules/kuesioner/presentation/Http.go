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
	BankSoaldomain "UnpakSiamida/modules/banksoal/domain"
	Kuesionerdomain "UnpakSiamida/modules/kuesioner/domain"

	ActiveKuesioners "UnpakSiamida/modules/kuesioner/application/ActiveKuesioner"
	ActiveKuesionersSingle "UnpakSiamida/modules/kuesioner/application/ActiveKuesionerSingle"
	CreateKuesioner "UnpakSiamida/modules/kuesioner/application/CreateKuesioner"
	DeleteKuesioner "UnpakSiamida/modules/kuesioner/application/DeleteKuesioner"
	GetAllKuesionerReport "UnpakSiamida/modules/kuesioner/application/GetAllKuesionerReport"
	GetAllKuesioners "UnpakSiamida/modules/kuesioner/application/GetAllKuesioners"
	GetKuesioner "UnpakSiamida/modules/kuesioner/application/GetKuesioner"
	GetKuesionerJawaban "UnpakSiamida/modules/kuesioner/application/GetKuesionerJawaban"
	SaveKuesionerJawaban "UnpakSiamida/modules/kuesioner/application/SaveKuesionerJawaban"
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
// POST /kuesioner/{uuid}/jawaban
// =======================================================

// SaveKuesionerJawabanHandler godoc
// @Summary Save new Kuesioner Jawaban
// @Tags Kuesioner
// @param pertanyaan formData string true "Pertanyaan" format(uuid)
// @param jawaban formData string true "Jawaban"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created Kuesioner"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /kuesioner [post]
func SaveKuesionerJawabanHandlerfunc(c *fiber.Ctx) error {

	UuidKuesioner := c.Params("uuid")
	UuidPertanyaan := c.FormValue("pertanyaan")
	Jawaban := c.FormValue("jawaban")

	SID := c.FormValue("sid")
	Resource := c.FormValue("resource")
	CodeCtx := c.FormValue("codectx")

	cmd := SaveKuesionerJawaban.SaveKuesionerJawabanCommand{
		UuidKuesioner:  UuidKuesioner,
		UuidPertanyaan: UuidPertanyaan,
		Jawaban:        Jawaban,
		SID:            SID,
		Resource:       Resource,
		CodeCtx:        CodeCtx,
	}

	uuid, err := mediatr.Send[SaveKuesionerJawaban.SaveKuesionerJawabanCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// GET /Kuesioner/{uuid}/jawaban
// =======================================================

// GetKuesionerJawabanHandler godoc
// @Summary Get Kuesioner Jawaban by UUID
// @Tags Kuesioner
// @Param uuid path string true "Kuesioner UUID" format(uuid)
// @Produce json
// @Success 200 {object} Kuesionerdomain.Kuesioner
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /Kuesioner/{uuid} [get]
func GetKuesionerJawabanHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetKuesionerJawaban.GetKuesionerJawabanByUuidQuery{Uuid: uuid}

	KuesionerJawaban, err := mediatr.Send[
		GetKuesionerJawaban.GetKuesionerJawabanByUuidQuery,
		[]Kuesionerdomain.KuesionerJawabanDefault,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(KuesionerJawaban)
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
// GET /kuesioners
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
// @Router /kuesioners [get]
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

// =======================================================
// POST /kuesioners/Report
// =======================================================

// GetAllKuesionersReportHandler godoc
// @Summary Get all Kuesioners For Report
// @Tags Kuesioner
// @param judul formData string true "Judul"
// @param semester formData string true "semster"
// @param is4year formData string true "is4year"
// @Produce json
// @Success 200 {object} commondomain.Paged[Kuesionerdomain.KuesionerDefault]
// @Router /kuesioners [post]
func GetAllKuesionersReportHandlerfunc(c *fiber.Ctx) error { //langsung sse
	judul := c.FormValue("judul")
	semester := c.FormValue("semester")
	is4year := c.FormValue("is4year") == "1"

	query := GetAllKuesionerReport.GetAllKuesionersReportQuery{
		JudulBankSoal: helper.StrPtr(judul),
		Semester:      helper.StrPtr(semester),
		Is4Year:       is4year,
	}

	result, err := mediatr.Send[
		GetAllKuesionerReport.GetAllKuesionersReportQuery,
		[]Kuesionerdomain.KuesionerResult,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	paged := commondomain.Paged[Kuesionerdomain.KuesionerResult]{
		Data: result,
	}

	var adapter commonpresentation.SSEAdapter[Kuesionerdomain.KuesionerResult]

	return adapter.Send(c, paged)
}

// =======================================================
// GET /kuesioners/Active
// =======================================================

// ActiveKuesionersHandler godoc
// @Summary Get all Kuesioners
// @Tags Kuesioner
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Kuesionerdomain.KuesionerDefault]
// @Router /kuesioners [get]
func ActiveKuesionersHandlerfunc(c *fiber.Ctx) error {
	// flag := c.Query("flag", "none") //with deleted
	nidn := c.FormValue("nidn")
	nip := c.FormValue("nip")
	npm := c.FormValue("npm")
	fakultas := c.FormValue("fakultas")
	prodi := c.FormValue("prodi")
	unit := strings.TrimSpace(c.FormValue("unit"))
	if strings.HasPrefix(strings.ToUpper(unit), "F.") || strings.HasPrefix(strings.ToUpper(unit), "F. ") {
		unit = "Fakultas " + strings.TrimSpace(unit[2:])
	}
	unit = strings.ToLower(unit)

	// withDeleted := false
	// if flag == "deleted" {
	// 	withDeleted = true
	// }
	query := ActiveKuesioners.ActiveKuesionerQuery{
		NPM:      helper.StrPtr(npm),
		NIDN:     helper.StrPtr(nidn),
		NIP:      helper.StrPtr(nip),
		Fakultas: helper.StrPtr(fakultas),
		Prodi:    helper.StrPtr(prodi),
		Unit:     helper.StrPtr(unit),
	}

	result, err := mediatr.Send[
		ActiveKuesioners.ActiveKuesionerQuery,
		[]BankSoaldomain.BankSoalDefault,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// GET /kuesioners/Active/{uuid}
// =======================================================

// GetKuesionersActiveByTargetHandler godoc
// @Summary Get Kuesioner Active
// @Tags Kuesioner
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[Kuesionerdomain.KuesionerDefault]
// @Router /kuesioners [get]
func GetKuesionersActiveByTargetHandler(c *fiber.Ctx) error {
	// flag := c.Query("flag", "none") //with deleted
	uuid := c.Params("uuid")
	nidn := c.FormValue("nidn")
	nip := c.FormValue("nip")
	npm := c.FormValue("npm")
	fakultas := c.FormValue("fakultas")
	prodi := c.FormValue("prodi")
	unit := strings.TrimSpace(c.FormValue("unit"))
	if strings.HasPrefix(strings.ToUpper(unit), "F.") || strings.HasPrefix(strings.ToUpper(unit), "F. ") {
		unit = "Fakultas " + strings.TrimSpace(unit[2:])
	}
	unit = strings.ToLower(unit)

	// withDeleted := false
	// if flag == "deleted" {
	// 	withDeleted = true
	// }
	query := ActiveKuesionersSingle.ActiveKuesionerSingleQuery{
		NPM:      helper.StrPtr(npm),
		NIDN:     helper.StrPtr(nidn),
		NIP:      helper.StrPtr(nip),
		Fakultas: helper.StrPtr(fakultas),
		Prodi:    helper.StrPtr(prodi),
		Unit:     helper.StrPtr(unit),
		UUID:     uuid,
	}

	result, err := mediatr.Send[
		ActiveKuesionersSingle.ActiveKuesionerSingleQuery,
		*BankSoaldomain.BankSoalDefault,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
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
	admin := []string{"admin"}
	whoamiURL := "http://127.0.0.1:3000/whoami"

	app.Get("/kuesioner/setupuuid", SetupUuidKuesionersHandlerfunc)

	app.Post("/kuesioner", commonpresentation.JWTMiddleware(), CreateKuesionerHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)
	app.Get("/kuesioner/:uuid/jawaban", commonpresentation.JWTMiddleware(), GetKuesionerJawabanHandlerfunc)
	app.Post("/kuesioner/:uuid/jawaban", commonpresentation.JWTMiddleware(), SaveKuesionerJawabanHandlerfunc)

	app.Delete("/kuesioner/:uuid", commonpresentation.JWTMiddleware(), DeleteKuesionerHandlerfunc)

	app.Get("/kuesioner/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetKuesionerHandlerfunc)
	app.Get("/kuesioners", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllKuesionersHandlerfunc)
	app.Post("/kuesioners/report", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllKuesionersReportHandlerfunc)
	app.Get("/kuesioners/active", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), ActiveKuesionersHandlerfunc)
	app.Get("/kuesioners/active/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), commonpresentation.RBACMiddleware(admin, whoamiURL), GetKuesionersActiveByTargetHandler)
}
