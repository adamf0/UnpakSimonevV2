package application

import (
	"context"

	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"
	"time"

	"github.com/google/uuid"
)

type DeleteKuesionerCommandHandler struct {
	Repo domainkuesioner.IKuesionerRepository
}

func (h *DeleteKuesionerCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteKuesionerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate UUID
	uuidKuesioner, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkuesioner.InvalidUuid()
	}

	// Delete by UUID
	if err := h.Repo.Delete(ctx, uuidKuesioner); err != nil { // FIXED
		return "", err
	}

	return cmd.Uuid, nil
}
