package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	copy "UnpakSiamida/modules/banksoal/application/CopyBankSoal"
	create "UnpakSiamida/modules/banksoal/application/CreateBankSoal"
	delete "UnpakSiamida/modules/banksoal/application/DeleteBankSoal"
	getAll "UnpakSiamida/modules/banksoal/application/GetAllBankSoals"
	get "UnpakSiamida/modules/banksoal/application/GetBankSoal"
	restore "UnpakSiamida/modules/banksoal/application/RestoreBankSoal"
	setupUuid "UnpakSiamida/modules/banksoal/application/SetupUuidBankSoal"
	update "UnpakSiamida/modules/banksoal/application/UpdateBankSoal"
	domainBankSoal "UnpakSiamida/modules/banksoal/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleBankSoal(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoBankSoal := NewBankSoalRepository(db)
	// if err := db.AutoMigrate(&domainBankSoal.BankSoal{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorBankSoal())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateBankSoalCommand,
		string,
	](&create.CreateBankSoalCommandHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateBankSoalCommand,
		string,
	](&update.UpdateBankSoalCommandHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		restore.RestoreBankSoalCommand,
		string,
	](&restore.RestoreBankSoalCommandHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		copy.CopyBankSoalCommand,
		string,
	](&copy.CopyBankSoalCommandHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteBankSoalCommand,
		string,
	](&delete.DeleteBankSoalCommandHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		get.GetBankSoalByUuidQuery,
		*domainBankSoal.BankSoal,
	](&get.GetBankSoalByUuidQueryHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllBankSoalsQuery,
		commondomain.Paged[domainBankSoal.BankSoalDefault],
	](&getAll.GetAllBankSoalsQueryHandler{
		Repo: repoBankSoal,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidBankSoalCommand,
		string,
	](&setupUuid.SetupUuidBankSoalCommandHandler{
		Repo: repoBankSoal,
	})

	commoninfra.RegisterValidation(create.CreateBankSoalCommandValidation, "BankSoalCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateBankSoalCommandValidation, "BankSoalUpdate.Validation")
	commoninfra.RegisterValidation(restore.RestoreBankSoalCommandValidation, "BankSoalRestore.Validation")
	commoninfra.RegisterValidation(copy.CopyBankSoalCommandValidation, "BankSoalCopy.Validation")
	commoninfra.RegisterValidation(delete.DeleteBankSoalCommandValidation, "BankSoalDelete.Validation")

	return nil
}
