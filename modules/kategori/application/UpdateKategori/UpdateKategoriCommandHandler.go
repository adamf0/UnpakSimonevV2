package application

import (
	"context"
	"errors"

	domainkategori "UnpakSiamida/modules/kategori/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateKategoriCommandHandler struct {
	Repo domainkategori.IKategoriRepository
}

func (h *UpdateKategoriCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateKategoriCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidKategori, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkategori.InvalidUuid()
	}

	var parentID *uint
	if cmd.SubKategori != nil {

		parentUUID, err := uuid.Parse(*cmd.SubKategori)
		if err != nil {
			return "", domainkategori.InvalidUuid()
		}

		parentKategori, err := h.Repo.GetByUuid(ctx, parentUUID)

		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", domainkategori.NotFound(*cmd.SubKategori)
			}

			return "", err
		}

		id := parentKategori.ID
		parentID = &id
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingKategori, err := h.Repo.GetByUuid(ctx, uuidKategori)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainkategori.NotFound(cmd.Uuid)
		}
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainkategori.UpdateKategori(
		existingKategori,
		uuidKategori,
		cmd.NamaKategori,
		parentID,
		cmd.Resource,
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedKategori := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	err = h.Repo.WithTx(ctx, func(repo domainkategori.IKategoriRepository) error {

		if err := repo.Update(ctx, updatedKategori); err != nil {
			return err
		}

		if err := repo.RebuildFullText(ctx); err != nil {
			return err
		}

		return nil
	})

	return updatedKategori.UUID.String(), nil
}
