package application

import (
	"context"
	"errors"

	"UnpakSiamida/common/helper"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateKategoriCommandHandler struct {
	Repo domainkategori.IKategoriRepository
}

func (h *CreateKategoriCommandHandler) Handle(
	ctx context.Context,
	cmd CreateKategoriCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var parentID *uint
	if cmd.SubKategori != nil {

		parentUUID, err := uuid.Parse(*cmd.SubKategori)
		if err != nil {
			return "", domainkategori.InvalidUuid()
		}

		existingKategori, err := h.Repo.GetByUuid(ctx, parentUUID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", domainkategori.NotFound(helper.NullableString(cmd.SubKategori))
			}
			return "", err
		}

		id := existingKategori.ID
		parentID = &id
	}

	result := domainkategori.NewKategori(
		cmd.NamaKategori,
		parentID,
		cmd.Resource,
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createKategori := result.Value

	err := h.Repo.WithTx(ctx, func(repo domainkategori.IKategoriRepository) error {

		if err := repo.Create(ctx, createKategori); err != nil {
			return err
		}

		if err := repo.RebuildFullText(ctx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
