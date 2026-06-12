package domaintest

import (
	"UnpakSiamida/common/helper"
	"UnpakSiamida/modules/account/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccount_TableName(t *testing.T) {
	acc := domain.Account{}
	assert.Equal(t, "users", acc.TableName())
}

func TestAccount_New(t *testing.T) {
	email := "test@unpak.ac.id"
	fakultas := "FKIP"
	prodi := "Pendidikan IPA"

	res := domain.NewAccount(
		"testuser",
		"password123",
		"user",
		"Test User",
		&email,
		&fakultas,
		&prodi,
	)

	assert.True(t, res.IsSuccess)
	assert.NotNil(t, res.Value)
	assert.Equal(t, "testuser", *res.Value.Username)
	assert.Equal(t, "password123", *res.Value.Password)
	assert.Equal(t, "user", *res.Value.Level)
	assert.Equal(t, "Test User", *res.Value.Name)
	assert.Equal(t, email, *res.Value.Email)
}

func TestAccount_Update(t *testing.T) {
	email := "test@unpak.ac.id"
	acc := domain.NewAccount("testuser", "password123", "user", "Test User", &email, nil, nil).Value

	// Success case with new password
	newPass := "newpassword"
	res := domain.UpdateAccount(
		acc,
		acc.UUID,
		"updateduser",
		&newPass,
		"admin",
		"Updated Name",
		&email,
		helper.StrPtr("Fakultas1"),
		helper.StrPtr("Prodi1"),
	)
	assert.True(t, res.IsSuccess)
	assert.Equal(t, "updateduser", *res.Value.Username)
	assert.Equal(t, "newpassword", *res.Value.Password)
	assert.Equal(t, "admin", *res.Value.Level)

	// Update with empty/nil password shouldn't overwrite existing
	resNoPass := domain.UpdateAccount(
		acc,
		acc.UUID,
		"updateduser",
		nil,
		"admin",
		"Updated Name",
		&email,
		nil,
		nil,
	)
	assert.True(t, resNoPass.IsSuccess)
	assert.Equal(t, "newpassword", *resNoPass.Value.Password)

	// Failure case: Nil object
	resNil := domain.UpdateAccount(
		nil,
		acc.UUID,
		"updateduser",
		nil,
		"admin",
		"Updated Name",
		&email,
		nil,
		nil,
	)
	assert.False(t, resNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resNil.Error)

	// Failure case: Mismatch UUID
	resMismatch := domain.UpdateAccount(
		acc,
		uuid.New(),
		"updateduser",
		nil,
		"admin",
		"Updated Name",
		&email,
		nil,
		nil,
	)
	assert.False(t, resMismatch.IsSuccess)
	assert.Equal(t, domain.InvalidData(), resMismatch.Error)
}

func TestAccount_DeleteAndRestore(t *testing.T) {
	acc := domain.NewAccount("testuser", "password123", "user", "Test User", nil, nil, nil).Value

	// Delete
	resDel := domain.DeleteAccount(acc)
	assert.True(t, resDel.IsSuccess)
	assert.NotNil(t, resDel.Value.DeletedAt)

	// Restore
	resRestore := domain.RestoreAccount(acc)
	assert.True(t, resRestore.IsSuccess)
	assert.Nil(t, resRestore.Value.DeletedAt)

	// Failure case: Nil object delete
	resDelNil := domain.DeleteAccount(nil)
	assert.False(t, resDelNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resDelNil.Error)

	// Failure case: Nil object restore
	resResNil := domain.RestoreAccount(nil)
	assert.False(t, resResNil.IsSuccess)
	assert.Equal(t, domain.EmptyData(), resResNil.Error)
}
