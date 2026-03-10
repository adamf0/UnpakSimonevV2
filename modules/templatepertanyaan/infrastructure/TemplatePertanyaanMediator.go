package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	copy "UnpakSiamida/modules/templatepertanyaan/application/CopyTemplatePertanyaan"
	create "UnpakSiamida/modules/templatepertanyaan/application/CreateTemplatePertanyaan"
	delete "UnpakSiamida/modules/templatepertanyaan/application/DeleteTemplatePertanyaan"
	getAll "UnpakSiamida/modules/templatepertanyaan/application/GetAllTemplatePertanyaans"
	get "UnpakSiamida/modules/templatepertanyaan/application/GetTemplatePertanyaan"
	restore "UnpakSiamida/modules/templatepertanyaan/application/RestoreTemplatePertanyaan"
	setupUuid "UnpakSiamida/modules/templatepertanyaan/application/SetupUuidTemplatePertanyaan"
	update "UnpakSiamida/modules/templatepertanyaan/application/UpdateTemplatePertanyaan"
	domainTemplatePertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"

	infraBankSoal "UnpakSiamida/modules/banksoal/infrastructure"
	infraKategori "UnpakSiamida/modules/kategori/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleTemplatePertanyaan(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoTemplatePertanyaan := NewTemplatePertanyaanRepository(db)
	repoKategori := infraKategori.NewKategoriRepository(db)
	repoBankSoal := infraBankSoal.NewBankSoalRepository(db)
	// if err := db.AutoMigrate(&domainTemplatePertanyaan.TemplatePertanyaan{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorTemplatePertanyaan())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateTemplatePertanyaanCommand,
		string,
	](&create.CreateTemplatePertanyaanCommandHandler{
		Repo:         repoTemplatePertanyaan,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateTemplatePertanyaanCommand,
		string,
	](&update.UpdateTemplatePertanyaanCommandHandler{
		Repo:         repoTemplatePertanyaan,
		RepoKategori: repoKategori,
		RepoBankSoal: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		restore.RestoreTemplatePertanyaanCommand,
		string,
	](&restore.RestoreTemplatePertanyaanCommandHandler{
		Repo: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		copy.CopyTemplatePertanyaanCommand,
		string,
	](&copy.CopyTemplatePertanyaanCommandHandler{
		Repo: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteTemplatePertanyaanCommand,
		string,
	](&delete.DeleteTemplatePertanyaanCommandHandler{
		Repo: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		get.GetTemplatePertanyaanByUuidQuery,
		*domainTemplatePertanyaan.TemplatePertanyaan,
	](&get.GetTemplatePertanyaanByUuidQueryHandler{
		Repo: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllTemplatePertanyaansQuery,
		commondomain.Paged[domainTemplatePertanyaan.TemplatePertanyaanDefault],
	](&getAll.GetAllTemplatePertanyaansQueryHandler{
		Repo: repoTemplatePertanyaan,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidTemplatePertanyaanCommand,
		string,
	](&setupUuid.SetupUuidTemplatePertanyaanCommandHandler{
		Repo: repoTemplatePertanyaan,
	})

	commoninfra.RegisterValidation(create.CreateTemplatePertanyaanCommandValidation, "TemplatePertanyaanCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateTemplatePertanyaanCommandValidation, "TemplatePertanyaanUpdate.Validation")
	commoninfra.RegisterValidation(restore.RestoreTemplatePertanyaanCommandValidation, "TemplatePertanyaanRestore.Validation")
	commoninfra.RegisterValidation(copy.CopyTemplatePertanyaanCommandValidation, "TemplatePertanyaanCopy.Validation")
	commoninfra.RegisterValidation(delete.DeleteTemplatePertanyaanCommandValidation, "TemplatePertanyaanDelete.Validation")

	return nil
}
