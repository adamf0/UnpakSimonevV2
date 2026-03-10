package infrastructure

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	login "UnpakSiamida/modules/account/application/Login"
	who "UnpakSiamida/modules/account/application/Whoami"
	domain "UnpakSiamida/modules/account/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleAccount(db *gorm.DB, dbSimak *gorm.DB, dbSimpeg *gorm.DB) error {
	repoAccount := NewAccountRepository(db, dbSimak, dbSimpeg)

	mediatr.RegisterRequestHandler[
		who.WhoamiCommand,
		*domain.Account,
	](&who.WhoamiCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		login.LoginCommand,
		*domain.LoginResult,
	](&login.LoginCommandHandler{
		Repo: repoAccount,
	})

	commoninfra.RegisterValidation(login.LoginCommandValidation, "Login.Validation")

	return nil
}
