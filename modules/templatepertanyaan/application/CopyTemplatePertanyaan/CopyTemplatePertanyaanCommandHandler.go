package application

import (
	"context"
	"errors"

	copyjawaban "UnpakSiamida/modules/templatejawaban/application/CopyTemplateJawaban"
	domaintemplatepertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"
	"time"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type CopyTemplatePertanyaanCommandHandler struct {
	Repo domaintemplatepertanyaan.ITemplatePertanyaanRepository
}

func (h *CopyTemplatePertanyaanCommandHandler) Handle(
	ctx context.Context,
	cmd CopyTemplatePertanyaanCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidTemplatePertanyaan, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplatepertanyaan.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingTemplatePertanyaan, err := h.Repo.GetByUuid(ctx, uuidTemplatePertanyaan)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaintemplatepertanyaan.NotFound(cmd.Uuid)
		}
		return "", err
	}
	copyCount, err := h.Repo.CountCopy(ctx, existingTemplatePertanyaan.Pertanyaan)
	if err != nil {
		return "", err
	}

	result := domaintemplatepertanyaan.CopyTemplatePertanyaan(
		existingTemplatePertanyaan,
		copyCount,
		cmd.Resource, //local, simak, simpeg
		cmd.SID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

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

	createTemplatePertanyaan := result.Value
	if err := repo.Create(ctx, createTemplatePertanyaan); err != nil {
		return "", err
	}

	cmdJawaban := copyjawaban.CopyTemplateJawabanCommand{
		Tx:                         tx,
		SourceTemplatePertanyaanID: existingTemplatePertanyaan.ID,
		TargetTemplatePertanyaanID: createTemplatePertanyaan.ID,
	}

	_, err = mediatr.Send[
		copyjawaban.CopyTemplateJawabanCommand,
		string,
	](
		ctx,
		cmdJawaban,
	)

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
