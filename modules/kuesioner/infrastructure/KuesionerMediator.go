package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	create "UnpakSiamida/modules/kuesioner/application/CreateKuesioner"
	delete "UnpakSiamida/modules/kuesioner/application/DeleteKuesioner"
	getAll "UnpakSiamida/modules/kuesioner/application/GetAllKuesioners"
	get "UnpakSiamida/modules/kuesioner/application/GetKuesioner"
	setupUuid "UnpakSiamida/modules/kuesioner/application/SetupUuidKuesioner"
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"

	infraAccount "UnpakSiamida/modules/account/infrastructure"
	infraBankSoal "UnpakSiamida/modules/banksoal/infrastructure"

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
	repoBankSoal := infraBankSoal.NewBankSoalRepository(db)
	repoAccount := infraAccount.NewAccountRepository(db, dbSimak, dbSimpeg)
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
		setupUuid.SetupUuidKuesionerCommand,
		string,
	](&setupUuid.SetupUuidKuesionerCommandHandler{
		Repo: repoKuesioner,
	})

	commoninfra.RegisterValidation(create.CreateKuesionerCommandValidation, "KuesionerCreate.Validation")
	commoninfra.RegisterValidation(delete.DeleteKuesionerCommandValidation, "KuesionerDelete.Validation")

	return nil
}
