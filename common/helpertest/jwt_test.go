package commontest

import (
	"UnpakSiamida/common/helper"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	sid := "user-123"
	resource := "simak"
	codectx := "dosen"

	accessToken, refreshToken, err := helper.GenerateToken(sid, resource, &codectx)
	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Validate the access token structure
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, sid, claims["sid"])
	assert.Equal(t, resource, claims["resource"])
	assert.Equal(t, codectx, claims["codectx"])
	assert.NotEmpty(t, claims["exp"])
}
