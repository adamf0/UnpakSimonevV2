package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	create "UnpakSiamida/modules/templatejawaban/application/CreateTemplateJawaban"
	delete "UnpakSiamida/modules/templatejawaban/application/DeleteTemplateJawaban"
	getAll "UnpakSiamida/modules/templatejawaban/application/GetAllTemplateJawabans"
	get "UnpakSiamida/modules/templatejawaban/application/GetTemplateJawaban"
	restore "UnpakSiamida/modules/templatejawaban/application/RestoreTemplateJawaban"
	setupUuid "UnpakSiamida/modules/templatejawaban/application/SetupUuidTemplateJawaban"
	update "UnpakSiamida/modules/templatejawaban/application/UpdateTemplateJawaban"
	domainTemplateJawaban "UnpakSiamida/modules/templatejawaban/domain"

	infraTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleTemplateJawaban(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoTemplateJawaban := NewTemplateJawabanRepository(db)
	repoTemplatePertanyaan := infraTemplatePertanyaan.NewTemplatePertanyaanRepository(db)
	// if err := db.AutoMigrate(&domainTemplateJawaban.TemplateJawaban{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorTemplateJawaban())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateTemplateJawabanCommand,
		string,
	](&create.CreateTemplateJawabanCommandHandler{
		Repo:           repoTemplateJawaban,
		RepoPertanyaan: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateTemplateJawabanCommand,
		string,
	](&update.UpdateTemplateJawabanCommandHandler{
		Repo:           repoTemplateJawaban,
		RepoPertanyaan: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		restore.RestoreTemplateJawabanCommand,
		string,
	](&restore.RestoreTemplateJawabanCommandHandler{
		Repo: repoTemplateJawaban,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteTemplateJawabanCommand,
		string,
	](&delete.DeleteTemplateJawabanCommandHandler{
		Repo: repoTemplateJawaban,
	})

	mediatr.RegisterRequestHandler[
		get.GetTemplateJawabanByUuidQuery,
		*domainTemplateJawaban.TemplateJawaban,
	](&get.GetTemplateJawabanByUuidQueryHandler{
		Repo: repoTemplateJawaban,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllTemplateJawabansQuery,
		commondomain.Paged[domainTemplateJawaban.TemplateJawabanDefault],
	](&getAll.GetAllTemplateJawabansQueryHandler{
		Repo: repoTemplateJawaban,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidTemplateJawabanCommand,
		string,
	](&setupUuid.SetupUuidTemplateJawabanCommandHandler{
		Repo: repoTemplateJawaban,
	})

	commoninfra.RegisterValidation(create.CreateTemplateJawabanCommandValidation, "TemplateJawabanCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateTemplateJawabanCommandValidation, "TemplateJawabanUpdate.Validation")
	commoninfra.RegisterValidation(restore.RestoreTemplateJawabanCommandValidation, "TemplateJawabanRestore.Validation")
	commoninfra.RegisterValidation(delete.DeleteTemplateJawabanCommandValidation, "TemplateJawabanDelete.Validation")

	return nil
}
