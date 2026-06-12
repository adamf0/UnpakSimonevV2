package infrastructure_test

import (
	"UnpakSiamida/common/infrastructure"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeTelegramClient(t *testing.T) {
	client := infrastructure.NewFakeTelegramClient()
	require.NotNil(t, client)

	err := client.SendHTML("Hello HTML")
	require.NoError(t, err)

	fake, ok := client.(*infrastructure.FakeTelegramClient)
	require.True(t, ok)
	assert.Len(t, fake.Messages, 1)
	assert.Equal(t, "Hello HTML", fake.Messages[0])
}

func TestDefaultTelegramFactory_Fake(t *testing.T) {
	factory := &infrastructure.DefaultTelegramFactory{UseFake: true}
	client, err := factory.Create()
	require.NoError(t, err)
	require.NotNil(t, client)

	_, ok := client.(*infrastructure.FakeTelegramClient)
	assert.True(t, ok)
}

func TestDefaultTelegramFactory_Real_NotConfigured(t *testing.T) {
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_CHAT_ID")

	factory := &infrastructure.DefaultTelegramFactory{UseFake: false}
	_, err := factory.Create()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "telegram env not configured")
}

func TestDefaultTelegramFactory_Real_Configured(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "dummy-token")
	os.Setenv("TELEGRAM_CHAT_ID", "dummy-chat-id")
	defer func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_CHAT_ID")
	}()

	factory := &infrastructure.DefaultTelegramFactory{UseFake: false}
	client, err := factory.Create()
	require.NoError(t, err)
	require.NotNil(t, client)

	// Sending will try to hit real Telegram endpoint and fail due to dummy token
	errSend := client.SendHTML("Hello")
	assert.Error(t, errSend)
}
