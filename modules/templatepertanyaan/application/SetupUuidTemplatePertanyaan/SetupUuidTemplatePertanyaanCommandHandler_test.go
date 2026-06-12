package application_test

import (
	"context"
	"errors"
	"testing"

	"UnpakSiamida/modules/templatepertanyaan/application/SetupUuidTemplatePertanyaan"
	"UnpakSiamida/modules/templatepertanyaan/application/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetupUuidTemplatePertanyaanCommandHandler_Success(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	repo.SetupUuidFunc = func(ctx context.Context) error {
		return nil
	}

	handler := &application.SetupUuidTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.SetupUuidTemplatePertanyaanCommand{}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "berhasil setup uuid pada data", res)
}

func TestSetupUuidTemplatePertanyaanCommandHandler_Fail(t *testing.T) {
	repo := &mock.MockTemplatePertanyaanRepository{}

	repo.SetupUuidFunc = func(ctx context.Context) error {
		return errors.New("setup error")
	}

	handler := &application.SetupUuidTemplatePertanyaanCommandHandler{
		Repo: repo,
	}

	cmd := application.SetupUuidTemplatePertanyaanCommand{}

	res, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "setup error")
	assert.Empty(t, res)
}
