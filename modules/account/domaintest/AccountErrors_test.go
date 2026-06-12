package domaintest

import (
	"testing"

	"UnpakSiamida/modules/account/domain"

	"github.com/stretchr/testify/assert"
)

func TestAccountErrors(t *testing.T) {
	assert.Equal(t, "Account.InvalidCredential", domain.InvalidCredential().Code)
	assert.Equal(t, "Account.NotFound", domain.NotFound("123").Code)
	assert.Equal(t, "Account.EmptyData", domain.EmptyData().Code)
	assert.Equal(t, "Account.InvalidData", domain.InvalidData().Code)
	assert.Equal(t, "Account.InvalidUuid", domain.InvalidUuid().Code)
}
