package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	login "UnpakSiamida/modules/account/application/Login"
	who "UnpakSiamida/modules/account/application/Whoami"
	domain "UnpakSiamida/modules/account/domain"

	create "UnpakSiamida/modules/account/application/CreateAccount"
	delete "UnpakSiamida/modules/account/application/DeleteAccount"
	get "UnpakSiamida/modules/account/application/GetAccount"
	getAll "UnpakSiamida/modules/account/application/GetAllAccounts"
	restore "UnpakSiamida/modules/account/application/RestoreAccount"
	setupUuid "UnpakSiamida/modules/account/application/SetupUuidAccount"
	update "UnpakSiamida/modules/account/application/UpdateAccount"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleAccount(db *gorm.DB, dbSimak *gorm.DB, dbSimpeg *gorm.DB) error {
	repoAccount := NewAccountRepository(db, dbSimak, dbSimpeg)

	mediatr.RegisterRequestHandler[
		who.WhoamiCommand,
		*domain.AccountDefault,
	](&who.WhoamiCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		login.LoginCommand,
		*domain.LoginResult,
	](&login.LoginCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		create.CreateAccountCommand,
		string,
	](&create.CreateAccountCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateAccountCommand,
		string,
	](&update.UpdateAccountCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		restore.RestoreAccountCommand,
		string,
	](&restore.RestoreAccountCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteAccountCommand,
		string,
	](&delete.DeleteAccountCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		get.GetAccountQuery,
		*domain.Account,
	](&get.GetAccountQueryHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllAccountsQuery,
		commondomain.Paged[domain.Account],
	](&getAll.GetAllAccountsQueryHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidAccountCommand,
		string,
	](&setupUuid.SetupUuidAccountCommandHandler{
		Repo: repoAccount,
	})

	commoninfra.RegisterValidation(login.LoginCommandValidation, "Login.Validation")
	commoninfra.RegisterValidation(create.CreateAccountCommandValidation, "AccountCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateAccountCommandValidation, "AccountUpdate.Validation")
	commoninfra.RegisterValidation(restore.RestoreAccountCommandValidation, "AccountRestore.Validation")
	commoninfra.RegisterValidation(delete.DeleteAccountCommandValidation, "AccountDelete.Validation")

	return nil
}
