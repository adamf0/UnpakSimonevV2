package presentation

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdihadeli/go-mediatr"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	login "UnpakSiamida/modules/account/application/Login"
	who "UnpakSiamida/modules/account/application/Whoami"
	domainaccount "UnpakSiamida/modules/account/domain"

	CreateAccount "UnpakSiamida/modules/account/application/CreateAccount"
	DeleteAccount "UnpakSiamida/modules/account/application/DeleteAccount"
	GetAccount "UnpakSiamida/modules/account/application/GetAccount"
	GetAllAccounts "UnpakSiamida/modules/account/application/GetAllAccounts"
	GetAllAccountsLdap "UnpakSiamida/modules/account/application/GetAllAccountsLdap"
	RestoreAccount "UnpakSiamida/modules/account/application/RestoreAccount"
	SetupUuidAccount "UnpakSiamida/modules/account/application/SetupUuidAccount"
	UpdateAccount "UnpakSiamida/modules/account/application/UpdateAccount"
)

// =======================================================
// POST /login
// =======================================================

// LoginHandler godoc
// @Summary Login
// @Tags Login
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Produce json
// @Success 200 {object} map[string]string "jwt"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /login [post]
func LoginHandlerfunc(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	cmd := login.LoginCommand{
		Username: username,
		Password: password,
	}

	result, err := mediatr.Send[login.LoginCommand, *domainaccount.LoginResult](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{
		// "user_id":       result.UserID,
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

func WhoAmIHandler(c *fiber.Ctx) error {
	userID := c.FormValue("sid")
	if userID == "" {
		userID = c.Query("sid")
	}
	if userID == "" {
		userID = c.FormValue("user_id")
	}
	if userID == "" {
		userID = c.Query("user_id")
	}
	if userID == "" {
		userID = string(c.Request().PostArgs().Peek("sid"))
	}
	if userID == "" {
		if s, ok := c.Locals("sid").(string); ok {
			userID = s
		}
	}

	resource := c.FormValue("resource")
	if resource == "" {
		resource = c.Query("resource")
	}
	if resource == "" {
		resource = string(c.Request().PostArgs().Peek("resource"))
	}
	if resource == "" {
		if r, ok := c.Locals("resource").(string); ok {
			resource = r
		}
	}

	codectx := c.FormValue("codectx")
	if codectx == "" {
		codectx = c.FormValue("codetx")
	}
	if codectx == "" {
		codectx = c.Query("codectx")
	}
	if codectx == "" {
		codectx = c.Query("codetx")
	}
	if codectx == "" {
		codectx = string(c.Request().PostArgs().Peek("codectx"))
	}
	if codectx == "" {
		codectx = string(c.Request().PostArgs().Peek("codetx"))
	}
	if codectx == "" {
		if cx, ok := c.Locals("codectx").(string); ok {
			codectx = cx
		} else if cx2, ok := c.Locals("codetx").(string); ok {
			codectx = cx2
		}
	}

	// Extract claims directly from Authorization or ctxtoken JWT header if fields are missing
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		authHeader = c.Query("ctxtoken")
	}
	if authHeader != "" {
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)
		if tokenStr != "" {
			token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
			if err == nil && token != nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userID == "" {
						if s, ok := claims["sid"].(string); ok && s != "" {
							userID = s
						} else if pref, ok := claims["preferred_username"].(string); ok && pref != "" {
							userID = pref
						} else if emp, ok := claims["employeeid"].(string); ok && emp != "" {
							userID = emp
						} else if sub, ok := claims["sub"].(string); ok && sub != "" {
							userID = sub
						}
					}
					if resource == "" {
						if r, ok := claims["resource"].(string); ok && r != "" {
							resource = r
						}
					}
					if codectx == "" {
						if cx, ok := claims["codectx"].(string); ok && cx != "" {
							codectx = cx
						} else if cx2, ok := claims["codetx"].(string); ok && cx2 != "" {
							codectx = cx2
						}
					}
				}
			}
		}
	}

	var (
		SID  *string
		NIP  *string
		NIDN *string
		NPM  *string
	)

	isDosen := strings.EqualFold(codectx, domainaccount.CtxDosen) || strings.EqualFold(codectx, "dosen")
	isMahasiswa := strings.EqualFold(codectx, domainaccount.CtxMahasiswa) || strings.EqualFold(codectx, "mahasiswa")
	isTendik := strings.EqualFold(codectx, "tendik") || strings.EqualFold(codectx, "pegawai") || strings.EqualFold(resource, "simpeg")

	if isDosen && (strings.EqualFold(resource, "simak") || resource == "") {
		NIDN = helper.StrPtr(userID)
	} else if isMahasiswa && (strings.EqualFold(resource, "simak") || resource == "") {
		NPM = helper.StrPtr(userID)
	} else if isTendik {
		NIP = helper.StrPtr(userID)
	}

	if userID != "" {
		SID = helper.StrPtr(userID)
	}

	cmd := who.WhoamiCommand{
		SID:  SID,
		NIM:  NPM,
		NIDN: NIDN,
		NIP:  NIP,
	}
	result, err := mediatr.Send[who.WhoamiCommand, *domainaccount.AccountDefault](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

// =======================================================
// POST /account
// =======================================================

// CreateAccountHandler godoc
// @Summary Create new Account
// @Tags Account
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param level formData string true "Level"
// @Param name formData string true "Name"
// @Param email formData string false "Email"
// @Param fakultas formData string false "Kode Fakultas"
// @Param prodi formData string false "Kode Prodi"
// @Produce json
// @Success 200 {object} map[string]string "uuid of created Account"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account [post]
func CreateAccountHandlerfunc(c *fiber.Ctx) error {

	username := c.FormValue("username")
	password := c.FormValue("password")
	level := c.FormValue("level")
	name := c.FormValue("name")

	var email *string
	var employeeID *string
	var fakultas *string
	var prodi *string

	if v := c.FormValue("email"); v != "" {
		email = helper.StrPtr(v)
	}

	if v := c.FormValue("employee_id"); v != "" {
		employeeID = helper.StrPtr(v)
	}

	if v := c.FormValue("fakultas"); v != "" {
		fakultas = helper.StrPtr(v)
	}

	if v := c.FormValue("prodi"); v != "" {
		prodi = helper.StrPtr(v)
	}

	cmd := CreateAccount.CreateAccountCommand{
		Username:   username,
		Password:   password,
		Level:      level,
		Name:       name,
		Email:      email,
		EmployeeID: employeeID,
		Fakultas:   fakultas,
		Prodi:      prodi,
	}

	uuid, err := mediatr.Send[
		CreateAccount.CreateAccountCommand,
		string,
	](context.Background(), cmd)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return commonpresentation.JsonUUID(c, uuid)
}

// =======================================================
// PUT /account/{uuid}
// =======================================================

// UpdateAccountHandler godoc
// @Summary Update existing Account
// @Tags Account
// @Param uuid path string true "Account UUID" format(uuid)
// @Param username formData string true "Username"
// @Param password formData string false "Password"
// @Param level formData string true "Level"
// @Param name formData string true "Name"
// @Param email formData string false "Email"
// @Param employee_id formData string false "Employee ID / NIP LDAP"
// @Param fakultas formData string false "Kode Fakultas"
// @Param prodi formData string false "Kode Prodi"
// @Produce json
// @Success 200 {object} map[string]string "uuid of updated Account"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account/{uuid} [put]
func UpdateAccountHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")
	username := c.FormValue("username")
	level := c.FormValue("level")
	name := c.FormValue("name")

	var password *string
	var email *string
	var employeeID *string
	var fakultas *string
	var prodi *string

	if v := c.FormValue("password"); v != "" {
		password = helper.StrPtr(v)
	}

	if v := c.FormValue("email"); v != "" {
		email = helper.StrPtr(v)
	}

	if v := c.FormValue("employee_id"); v != "" {
		employeeID = helper.StrPtr(v)
	}

	if v := c.FormValue("fakultas"); v != "" {
		fakultas = helper.StrPtr(v)
	}

	if v := c.FormValue("prodi"); v != "" {
		prodi = helper.StrPtr(v)
	}

	cmd := UpdateAccount.UpdateAccountCommand{
		Uuid:       uuid,
		Username:   username,
		Password:   password,
		Level:      level,
		Name:       name,
		Email:      email,
		EmployeeID: employeeID,
		Fakultas:   fakultas,
		Prodi:      prodi,
	}

	updatedID, err := mediatr.Send[UpdateAccount.UpdateAccountCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": updatedID})
}

// =======================================================
// DELETE /account/{uuid}
// =======================================================

// DeleteAccountHandler godoc
// @Summary Delete a Account
// @Tags Account
// @Param uuid path string true "Account UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Account"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account/{uuid} [delete]
func DeleteAccountHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteAccount.DeleteAccountCommand{
		Uuid: uuid,
		Mode: "soft_delete",
	}

	deletedID, err := mediatr.Send[DeleteAccount.DeleteAccountCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// DELETE /account/{uuid}/force
// =======================================================

// ForceAccountHandler godoc
// @Summary Delete a Account
// @Tags Account
// @Param uuid path string true "Account UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Account"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account/{uuid} [delete]
func ForceDeleteAccountHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := DeleteAccount.DeleteAccountCommand{
		Uuid: uuid,
		Mode: "hard_delete",
	}

	deletedID, err := mediatr.Send[DeleteAccount.DeleteAccountCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// PUT /account/{uuid}/restore
// =======================================================

// RestoreAccountHandler godoc
// @Summary Restore a Account
// @Tags Account
// @Param uuid path string true "Account UUID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "uuid of deleted Account"
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account/{uuid} [delete]
func RestoreAccountHandlerfunc(c *fiber.Ctx) error {

	uuid := c.Params("uuid")

	cmd := RestoreAccount.RestoreAccountCommand{
		Uuid: uuid,
	}

	deletedID, err := mediatr.Send[RestoreAccount.RestoreAccountCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"uuid": deletedID})
}

// =======================================================
// GET /account/{uuid}
// =======================================================

// GetAccountHandler godoc
// @Summary Get Account by UUID
// @Tags Account
// @Param uuid path string true "Account UUID" format(uuid)
// @Produce json
// @Success 200 {object} domainaccount.Account
// @Failure 400 {object} commoninfra.ResponseError
// @Failure 404 {object} commoninfra.ResponseError
// @Failure 409 {object} commoninfra.ResponseError
// @Failure 500 {object} commoninfra.ResponseError
// @Router /account/{uuid} [get]
func GetAccountHandlerfunc(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	query := GetAccount.GetAccountQuery{Uuid: uuid}

	Account, err := mediatr.Send[
		GetAccount.GetAccountQuery,
		*domainaccount.Account,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if Account == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
	}

	return c.JSON(Account)
}

// =======================================================
// GET /accounts
// =======================================================

// GetAllAccountsHandler godoc
// @Summary Get all Accounts
// @Tags Account
// @Param mode query string false "paging | all | ndjson | sse"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Produce json
// @Success 200 {object} commondomain.Paged[domainaccount.Account]
// @Router /accounts [get]
func GetAllAccountsHandlerfunc(c *fiber.Ctx) error {
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
	query := GetAllAccounts.GetAllAccountsQuery{
		Search:        search,
		SearchFilters: filters,
		Deleted:       withDeleted,
	}

	var adapter commonpresentation.OutputAdapter[domainaccount.Account]
	switch mode {
	case "all":
		adapter = &commonpresentation.AllAdapter[domainaccount.Account]{}
	case "ndjson":
		adapter = &commonpresentation.NDJSONAdapter[domainaccount.Account]{}
	case "sse":
		adapter = &commonpresentation.SSEAdapter[domainaccount.Account]{}
	default:
		query.Page = &page
		query.Limit = &limit
		adapter = &commonpresentation.PagingAdapter[domainaccount.Account]{}
	}

	result, err := mediatr.Send[
		GetAllAccounts.GetAllAccountsQuery,
		commondomain.Paged[domainaccount.Account],
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return adapter.Send(c, result)
}

func GetAllAccountsLdapHandlerfunc(c *fiber.Ctx) error {
	groupsParam := c.Query("groups", "")
	var groups []string
	if groupsParam != "" {
		parts := strings.Split(groupsParam, ",")
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				groups = append(groups, trimmed)
			}
		}
	}

	query := GetAllAccountsLdap.GetAllAccountsLdapQuery{
		Groups: groups,
	}

	result, err := mediatr.Send[
		GetAllAccountsLdap.GetAllAccountsLdapQuery,
		[]domainaccount.AccountLdap,
	](context.Background(), query)

	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	if result == nil {
		result = []domainaccount.AccountLdap{}
	}

	return c.JSON(fiber.Map{
		"data":  result,
		"total": len(result),
	})
}

func SetupUuidAccountsHandlerfunc(c *fiber.Ctx) error {
	cmd := SetupUuidAccount.SetupUuidAccountCommand{}

	message, err := mediatr.Send[SetupUuidAccount.SetupUuidAccountCommand, string](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(fiber.Map{"message": message})
}

func ModuleAccount(app *fiber.App) {
	app.Post("/api/v2/login", LoginHandlerfunc)
	app.Get("/api/v2/whoami", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), WhoAmIHandler)

	app.Get("/api/v2/account/setupuuid", SetupUuidAccountsHandlerfunc)

	app.Post("/api/v2/account", commonpresentation.JWTMiddleware(), CreateAccountHandlerfunc) //commonpresentation.RBACMiddleware(admin, whoamiURL)
	app.Put("/api/v2/account/:uuid", commonpresentation.JWTMiddleware(), UpdateAccountHandlerfunc)

	app.Delete("/api/v2/account/:uuid", commonpresentation.JWTMiddleware(), DeleteAccountHandlerfunc)            //soft delete
	app.Delete("/api/v2/account/:uuid/force", commonpresentation.JWTMiddleware(), ForceDeleteAccountHandlerfunc) //hanya lpm saja yg hard delete
	app.Put("/api/v2/account/:uuid/restore", commonpresentation.JWTMiddleware(), RestoreAccountHandlerfunc)

	app.Get("/api/v2/account/:uuid", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAccountHandlerfunc)
	app.Get("/api/v2/accounts", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), GetAllAccountsHandlerfunc)
	app.Get("/api/v2/accounts/ldap", commonpresentation.SmartCompress(), GetAllAccountsLdapHandlerfunc)
}
