package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	getAll "UnpakSiamida/modules/fakultas/application/GetAllFakultass"
	domainFakultas "UnpakSiamida/modules/fakultas/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleFakultas(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoFakultas := NewFakultasRepository(db)
	// if err := db.AutoMigrate(&domainFakultas.Fakultas{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorFakultas())

	// Register request handler
	mediatr.RegisterRequestHandler[
		getAll.GetAllFakultassQuery,
		commondomain.Paged[domainFakultas.FakultasDefault],
	](&getAll.GetAllFakultassQueryHandler{
		Repo: repoFakultas,
	})

	return nil
}
