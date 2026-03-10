package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	copy "UnpakSiamida/modules/kategori/application/CopyKategori"
	create "UnpakSiamida/modules/kategori/application/CreateKategori"
	delete "UnpakSiamida/modules/kategori/application/DeleteKategori"
	getAll "UnpakSiamida/modules/kategori/application/GetAllKategoris"
	get "UnpakSiamida/modules/kategori/application/GetKategori"
	restore "UnpakSiamida/modules/kategori/application/RestoreKategori"
	setupUuid "UnpakSiamida/modules/kategori/application/SetupUuidKategori"
	update "UnpakSiamida/modules/kategori/application/UpdateKategori"
	domainKategori "UnpakSiamida/modules/kategori/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleKategori(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoKategori := NewKategoriRepository(db)
	// if err := db.AutoMigrate(&domainKategori.Kategori{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorKategori())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateKategoriCommand,
		string,
	](&create.CreateKategoriCommandHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateKategoriCommand,
		string,
	](&update.UpdateKategoriCommandHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		restore.RestoreKategoriCommand,
		string,
	](&restore.RestoreKategoriCommandHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		copy.CopyKategoriCommand,
		string,
	](&copy.CopyKategoriCommandHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteKategoriCommand,
		string,
	](&delete.DeleteKategoriCommandHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		get.GetKategoriByUuidQuery,
		*domainKategori.Kategori,
	](&get.GetKategoriByUuidQueryHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllKategorisQuery,
		commondomain.Paged[domainKategori.KategoriDefault],
	](&getAll.GetAllKategorisQueryHandler{
		Repo: repoKategori,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidKategoriCommand,
		string,
	](&setupUuid.SetupUuidKategoriCommandHandler{
		Repo: repoKategori,
	})

	commoninfra.RegisterValidation(create.CreateKategoriCommandValidation, "KategoriCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateKategoriCommandValidation, "KategoriUpdate.Validation")
	commoninfra.RegisterValidation(restore.RestoreKategoriCommandValidation, "KategoriRestore.Validation")
	commoninfra.RegisterValidation(copy.CopyKategoriCommandValidation, "KategoriCopy.Validation")
	commoninfra.RegisterValidation(delete.DeleteKategoriCommandValidation, "KategoriDelete.Validation")

	return nil
}
