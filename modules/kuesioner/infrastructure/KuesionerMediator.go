package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"
	getActive "UnpakSiamida/modules/kuesioner/application/ActiveKuesioner"
	getActiveSingle "UnpakSiamida/modules/kuesioner/application/ActiveKuesionerSingle"
	create "UnpakSiamida/modules/kuesioner/application/CreateKuesioner"
	delete "UnpakSiamida/modules/kuesioner/application/DeleteKuesioner"
	getReport "UnpakSiamida/modules/kuesioner/application/GetAllKuesionerReport"
	getAll "UnpakSiamida/modules/kuesioner/application/GetAllKuesioners"
	get "UnpakSiamida/modules/kuesioner/application/GetKuesioner"
	getAnsware "UnpakSiamida/modules/kuesioner/application/GetKuesionerJawaban"
	save "UnpakSiamida/modules/kuesioner/application/SaveKuesionerJawaban"
	setupUuid "UnpakSiamida/modules/kuesioner/application/SetupUuidKuesioner"
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"

	infraAccount "UnpakSiamida/modules/account/infrastructure"
	infraBankSoal "UnpakSiamida/modules/banksoal/infrastructure"
	infraJawaban "UnpakSiamida/modules/templatejawaban/infrastructure"
	infraPertanyaan "UnpakSiamida/modules/templatepertanyaan/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleKuesioner(db *gorm.DB, dbSimak *gorm.DB, dbSimpeg *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoKuesioner := NewKuesionerRepository(db)
	repoKuesionerJawaban := NewKuesionerJawabanRepository(db)
	repoBankSoal := infraBankSoal.NewBankSoalRepository(db)
	repoAccount := infraAccount.NewAccountRepository(db, dbSimak, dbSimpeg)
	repoPertanyaan := infraPertanyaan.NewTemplatePertanyaanRepository(db)
	repoJawaban := infraJawaban.NewTemplateJawabanRepository(db)
	// if err := db.AutoMigrate(&domainKuesioner.Kuesioner{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorKuesioner())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateKuesionerCommand,
		string,
	](&create.CreateKuesionerCommandHandler{
		Repo:         repoKuesioner,
		RepoBankSoal: repoBankSoal,
		RepoAccount:  repoAccount,
	})

	mediatr.RegisterRequestHandler[
		save.SaveKuesionerJawabanCommand,
		string,
	](&save.SaveKuesionerJawabanCommandHandler{
		Repo:                 repoKuesioner,
		RepoPertanyaan:       repoPertanyaan,
		RepoJawaban:          repoJawaban,
		RepoJawabanKuesioner: repoKuesionerJawaban,
	})

	mediatr.RegisterRequestHandler[
		getReport.GetAllKuesionersReportQuery,
		[]domainKuesioner.KuesionerResult,
	](&getReport.GetAllKuesionersReportQueryHandler{
		Repo: repoKuesioner,
	})

	mediatr.RegisterRequestHandler[
		getAnsware.GetKuesionerJawabanByUuidQuery,
		[]domainKuesioner.KuesionerJawabanDefault,
	](&getAnsware.GetKuesionerJawabanByUuidQueryQueryHandler{
		Repo: repoKuesionerJawaban,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteKuesionerCommand,
		string,
	](&delete.DeleteKuesionerCommandHandler{
		Repo: repoKuesioner,
	})

	mediatr.RegisterRequestHandler[
		get.GetKuesionerByUuidQuery,
		*domainKuesioner.Kuesioner,
	](&get.GetKuesionerByUuidQueryHandler{
		Repo: repoKuesioner,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllKuesionersQuery,
		commondomain.Paged[domainKuesioner.KuesionerDefault],
	](&getAll.GetAllKuesionersQueryHandler{
		Repo: repoKuesioner,
	})

	mediatr.RegisterRequestHandler[
		getActive.ActiveKuesionerQuery,
		[]domainBankSoal.BankSoalDefault,
	](&getActive.ActiveKuesionerQueryHandler{
		Repo:         repoKuesioner,
		RepoJawaban:  repoKuesionerJawaban,
		RepoBankSoal: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		getActiveSingle.ActiveKuesionerSingleQuery,
		*domainBankSoal.BankSoalDefault,
	](&getActiveSingle.ActiveKuesionerSingleQueryHandler{
		Repo:         repoKuesioner,
		RepoJawaban:  repoKuesionerJawaban,
		RepoBankSoal: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidKuesionerCommand,
		string,
	](&setupUuid.SetupUuidKuesionerCommandHandler{
		Repo: repoKuesioner,
	})

	commoninfra.RegisterValidation(create.CreateKuesionerCommandValidation, "KuesionerCreate.Validation")
	commoninfra.RegisterValidation(save.SaveKuesionerJawabanCommandValidation, "KuesionerJawabanSave.Validation")
	commoninfra.RegisterValidation(delete.DeleteKuesionerCommandValidation, "KuesionerDelete.Validation")
	commoninfra.RegisterValidation(getReport.GetAllKuesionersReportQueryValidation, "KuesionerReport.Validation")

	return nil
}
