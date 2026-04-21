package presentation

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"

	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
	login "UnpakSiamida/modules/account/application/Login"
	who "UnpakSiamida/modules/account/application/Whoami"
	domainaccount "UnpakSiamida/modules/account/domain"
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
	resource := c.FormValue("resource")
	codectx := c.FormValue("codectx")

	var (
		SID  *string
		NIP  *string
		NIDN *string
		NPM  *string
	)

	if codectx == domainaccount.CtxDosen && helper.NullableString(&resource) == "simak" {
		NIDN = helper.StrPtr(userID)
	} else if codectx == domainaccount.CtxMahasiswa && helper.NullableString(&resource) == "simak" {
		NPM = helper.StrPtr(userID)
	} else if resource == "simpeg" {
		NIP = helper.StrPtr(userID)
	} else {
		SID = helper.StrPtr(userID)
	}

	cmd := who.WhoamiCommand{
		SID:  SID,
		NIM:  NPM,
		NIDN: NIDN,
		NIP:  NIP,
	}
	result, err := mediatr.Send[who.WhoamiCommand, *domainaccount.Account](context.Background(), cmd)
	if err != nil {
		return commoninfra.HandleError(c, err)
	}

	return c.JSON(result)
}

func ModuleAccount(app *fiber.App) {
	app.Post("/login", LoginHandlerfunc)
	app.Get("/whoami", commonpresentation.SmartCompress(), commonpresentation.JWTMiddleware(), WhoAmIHandler)
}
