package presentation_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	commondomain "UnpakSiamida/common/domain"
	login "UnpakSiamida/modules/account/application/Login"
	who "UnpakSiamida/modules/account/application/Whoami"
	"UnpakSiamida/modules/account/domain"
	"UnpakSiamida/modules/account/presentation"

	CreateAccount "UnpakSiamida/modules/account/application/CreateAccount"
	DeleteAccount "UnpakSiamida/modules/account/application/DeleteAccount"
	GetAccount "UnpakSiamida/modules/account/application/GetAccount"
	GetAllAccounts "UnpakSiamida/modules/account/application/GetAllAccounts"
	RestoreAccount "UnpakSiamida/modules/account/application/RestoreAccount"
	SetupUuidAccount "UnpakSiamida/modules/account/application/SetupUuidAccount"
	UpdateAccount "UnpakSiamida/modules/account/application/UpdateAccount"
)

var (
	mockLoginFunc            func(ctx context.Context, cmd login.LoginCommand) (*domain.LoginResult, error)
	mockWhoamiFunc           func(ctx context.Context, cmd who.WhoamiCommand) (*domain.AccountDefault, error)
	mockCreateAccountFunc    func(ctx context.Context, cmd CreateAccount.CreateAccountCommand) (string, error)
	mockUpdateAccountFunc    func(ctx context.Context, cmd UpdateAccount.UpdateAccountCommand) (string, error)
	mockDeleteAccountFunc    func(ctx context.Context, cmd DeleteAccount.DeleteAccountCommand) (string, error)
	mockRestoreAccountFunc   func(ctx context.Context, cmd RestoreAccount.RestoreAccountCommand) (string, error)
	mockGetAccountFunc       func(ctx context.Context, q GetAccount.GetAccountQuery) (*domain.Account, error)
	mockGetAllAccountsFunc   func(ctx context.Context, q GetAllAccounts.GetAllAccountsQuery) (commondomain.Paged[domain.Account], error)
	mockSetupUuidAccountFunc func(ctx context.Context, cmd SetupUuidAccount.SetupUuidAccountCommand) (string, error)

	registerOnce sync.Once
)

type mockLoginHandler struct{}

func (h *mockLoginHandler) Handle(ctx context.Context, cmd login.LoginCommand) (*domain.LoginResult, error) {
	return mockLoginFunc(ctx, cmd)
}

type mockWhoamiHandler struct{}

func (h *mockWhoamiHandler) Handle(ctx context.Context, cmd who.WhoamiCommand) (*domain.AccountDefault, error) {
	return mockWhoamiFunc(ctx, cmd)
}

type mockCreateAccountHandler struct{}

func (h *mockCreateAccountHandler) Handle(ctx context.Context, cmd CreateAccount.CreateAccountCommand) (string, error) {
	return mockCreateAccountFunc(ctx, cmd)
}

type mockUpdateAccountHandler struct{}

func (h *mockUpdateAccountHandler) Handle(ctx context.Context, cmd UpdateAccount.UpdateAccountCommand) (string, error) {
	return mockUpdateAccountFunc(ctx, cmd)
}

type mockDeleteAccountHandler struct{}

func (h *mockDeleteAccountHandler) Handle(ctx context.Context, cmd DeleteAccount.DeleteAccountCommand) (string, error) {
	return mockDeleteAccountFunc(ctx, cmd)
}

type mockRestoreAccountHandler struct{}

func (h *mockRestoreAccountHandler) Handle(ctx context.Context, cmd RestoreAccount.RestoreAccountCommand) (string, error) {
	return mockRestoreAccountFunc(ctx, cmd)
}

type mockGetAccountHandler struct{}

func (h *mockGetAccountHandler) Handle(ctx context.Context, q GetAccount.GetAccountQuery) (*domain.Account, error) {
	return mockGetAccountFunc(ctx, q)
}

type mockGetAllAccountsHandler struct{}

func (h *mockGetAllAccountsHandler) Handle(ctx context.Context, q GetAllAccounts.GetAllAccountsQuery) (commondomain.Paged[domain.Account], error) {
	return mockGetAllAccountsFunc(ctx, q)
}

type mockSetupUuidAccountHandler struct{}

func (h *mockSetupUuidAccountHandler) Handle(ctx context.Context, cmd SetupUuidAccount.SetupUuidAccountCommand) (string, error) {
	return mockSetupUuidAccountFunc(ctx, cmd)
}

func setupMediatrMocks() {
	registerOnce.Do(func() {
		_ = mediatr.RegisterRequestHandler[login.LoginCommand, *domain.LoginResult](&mockLoginHandler{})
		_ = mediatr.RegisterRequestHandler[who.WhoamiCommand, *domain.AccountDefault](&mockWhoamiHandler{})
		_ = mediatr.RegisterRequestHandler[CreateAccount.CreateAccountCommand, string](&mockCreateAccountHandler{})
		_ = mediatr.RegisterRequestHandler[UpdateAccount.UpdateAccountCommand, string](&mockUpdateAccountHandler{})
		_ = mediatr.RegisterRequestHandler[DeleteAccount.DeleteAccountCommand, string](&mockDeleteAccountHandler{})
		_ = mediatr.RegisterRequestHandler[RestoreAccount.RestoreAccountCommand, string](&mockRestoreAccountHandler{})
		_ = mediatr.RegisterRequestHandler[GetAccount.GetAccountQuery, *domain.Account](&mockGetAccountHandler{})
		_ = mediatr.RegisterRequestHandler[GetAllAccounts.GetAllAccountsQuery, commondomain.Paged[domain.Account]](&mockGetAllAccountsHandler{})
		_ = mediatr.RegisterRequestHandler[SetupUuidAccount.SetupUuidAccountCommand, string](&mockSetupUuidAccountHandler{})
	})
}

func generateToken(sid, resource, codectx string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":      sid,
		"resource": resource,
		"codectx":  codectx,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}

func TestAccountPresentation(t *testing.T) {
	setupMediatrMocks()

	app := fiber.New()
	presentation.ModuleAccount(app)

	token := generateToken("user-123", "simak", domain.CtxDosen)

	t.Run("Login success", func(t *testing.T) {
		mockLoginFunc = func(ctx context.Context, cmd login.LoginCommand) (*domain.LoginResult, error) {
			assert.Equal(t, "admin", cmd.Username)
			return &domain.LoginResult{
				AccessToken:  "token-123",
				RefreshToken: "refresh-123",
			}, nil
		}

		form := url.Values{}
		form.Add("username", "admin")
		form.Add("password", "pass")

		req := httptest.NewRequest("POST", "/api/v2/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "token-123", res["access_token"])
		assert.Equal(t, "refresh-123", res["refresh_token"])
	})

	t.Run("Whoami success", func(t *testing.T) {
		mockWhoamiFunc = func(ctx context.Context, cmd who.WhoamiCommand) (*domain.AccountDefault, error) {
			assert.Equal(t, "user-123", *cmd.NIDN)
			name := "John Doe"
			return &domain.AccountDefault{
				Name: &name,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/whoami", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.AccountDefault
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "John Doe", *res.Name)
	})

	t.Run("Setup UUID Account success", func(t *testing.T) {
		mockSetupUuidAccountFunc = func(ctx context.Context, cmd SetupUuidAccount.SetupUuidAccountCommand) (string, error) {
			return "setup complete", nil
		}

		req := httptest.NewRequest("GET", "/api/v2/account/setupuuid", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "setup complete", res["message"])
	})

	t.Run("Create Account success", func(t *testing.T) {
		mockCreateAccountFunc = func(ctx context.Context, cmd CreateAccount.CreateAccountCommand) (string, error) {
			assert.Equal(t, "new_user", cmd.Username)
			return "new-account-uuid", nil
		}

		form := url.Values{}
		form.Add("username", "new_user")
		form.Add("password", "pass")
		form.Add("level", "user")
		form.Add("name", "New User")

		req := httptest.NewRequest("POST", "/api/v2/account", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "new-account-uuid", res["uuid"])
	})

	t.Run("Update Account success", func(t *testing.T) {
		mockUpdateAccountFunc = func(ctx context.Context, cmd UpdateAccount.UpdateAccountCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "updated_user", cmd.Username)
			return "uuid-123", nil
		}

		form := url.Values{}
		form.Add("username", "updated_user")
		form.Add("level", "user")
		form.Add("name", "Updated User")

		req := httptest.NewRequest("PUT", "/api/v2/account/uuid-123", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete Account soft success", func(t *testing.T) {
		mockDeleteAccountFunc = func(ctx context.Context, cmd DeleteAccount.DeleteAccountCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "soft_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/account/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Delete Account hard success", func(t *testing.T) {
		mockDeleteAccountFunc = func(ctx context.Context, cmd DeleteAccount.DeleteAccountCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			assert.Equal(t, "hard_delete", cmd.Mode)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("DELETE", "/api/v2/account/uuid-123/force", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Restore Account success", func(t *testing.T) {
		mockRestoreAccountFunc = func(ctx context.Context, cmd RestoreAccount.RestoreAccountCommand) (string, error) {
			assert.Equal(t, "uuid-123", cmd.Uuid)
			return "uuid-123", nil
		}

		req := httptest.NewRequest("PUT", "/api/v2/account/uuid-123/restore", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "uuid-123", res["uuid"])
	})

	t.Run("Get Account by UUID success", func(t *testing.T) {
		username := "john"
		accountVal := &domain.Account{
			Username: &username,
		}
		mockGetAccountFunc = func(ctx context.Context, q GetAccount.GetAccountQuery) (*domain.Account, error) {
			assert.Equal(t, "uuid-123", q.Uuid)
			return accountVal, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/account/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res domain.Account
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "john", *res.Username)
	})

	t.Run("Get Account by UUID not found", func(t *testing.T) {
		mockGetAccountFunc = func(ctx context.Context, q GetAccount.GetAccountQuery) (*domain.Account, error) {
			return nil, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/account/uuid-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get All Accounts success", func(t *testing.T) {
		username1 := "user1"
		username2 := "user2"
		items := []domain.Account{
			{Username: &username1},
			{Username: &username2},
		}
		mockGetAllAccountsFunc = func(ctx context.Context, q GetAllAccounts.GetAllAccountsQuery) (commondomain.Paged[domain.Account], error) {
			return commondomain.Paged[domain.Account]{
				Data:  items,
				Total: 2,
			}, nil
		}

		req := httptest.NewRequest("GET", "/api/v2/accounts?mode=all", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []domain.Account
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res, 2)
		assert.Equal(t, "user1", *res[0].Username)
	})
}
