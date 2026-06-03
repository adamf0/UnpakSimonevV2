package application

import (
	"context"
	"errors"

	domainkategori "UnpakSiamida/modules/kategori/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CopyKategoriCommandHandler struct {
	Repo domainkategori.IKategoriRepository
}

func (h *CopyKategoriCommandHandler) Handle(
	ctx context.Context,
	cmd CopyKategoriCommand,
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
	copyCount, err := h.Repo.CountCopy(ctx, existingKategori.NamaKategori)
	if err != nil {
		return "", err
	}

	result := domainkategori.CopyKategori(
		existingKategori,
		copyCount,
		cmd.Resource, //local, simak, simpeg
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createKategori := result.Value

	// err = h.Repo.WithTx(ctx, func(repo domainkategori.IKategoriRepository) error {

	// 	if err := repo.Create(ctx, createKategori); err != nil {
	// 		return err
	// 	}

	// 	if err := repo.RebuildFullText(ctx); err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	tx, err := h.Repo.BeginTx(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	commit := false

	defer func() {
		if !commit {
			_ = tx.Rollback()
		}
	}()

	repo := h.Repo.WithTx(tx)

	if err := repo.Create(ctx, createKategori); err != nil {
		return "", err
	}

	if err := repo.RebuildFullText(ctx); err != nil {
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	commit = true

	return result.Value.UUID.String(), nil
}
