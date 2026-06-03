package application

import (
	"context"
	"errors"

	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	copyJawaban "UnpakSiamida/modules/templatejawaban/application/CopyTemplateJawaban"
	copy "UnpakSiamida/modules/templatepertanyaan/application/CopyTemplatePertanyaan"

	"time"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type CopyBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *CopyBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd CopyBankSoalCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidBankSoal, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainbanksoal.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	existingBankSoal, err := h.Repo.GetByUuid(ctx, uuidBankSoal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainbanksoal.NotFound(cmd.Uuid)
		}
		return "", err
	}
	copyCount, err := h.Repo.CountCopy(ctx, existingBankSoal.Judul)
	if err != nil {
		return "", err
	}

	result := domainbanksoal.CopyBankSoal(
		existingBankSoal,
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
	// repoPertanyaan := h.RepoPertanyaan.WithTx(tx)
	// repoJawaban := h.RepoJawaban.WithTx(tx)

	createBankSoal := result.Value
	if err := repo.Create(ctx, createBankSoal); err != nil {
		return "", err
	}

	// -------------------------
	// COPY PERTANYAAN
	// -------------------------

	mapping, err := mediatr.Send[
		copy.CopyTemplatePertanyaanResultCommand,
		map[uint]uint,
	](ctx, copy.CopyTemplatePertanyaanResultCommand{
		Tx:               tx,
		SourceBankSoalID: existingBankSoal.ID,
		TargetBankSoalID: createBankSoal.ID,
		Resource:         cmd.Resource,
		Sid:              cmd.SID,
	})
	if err != nil {
		return "", err
	}

	// -------------------------
	// COPY JAWABAN
	// -------------------------
	for oldPertanyaanID, newPertanyaanID := range mapping {

		cmdJawaban := copyJawaban.CopyTemplateJawabanCommand{
			Tx:                         tx,
			SourceTemplatePertanyaanID: oldPertanyaanID,
			TargetTemplatePertanyaanID: newPertanyaanID,
		}

		_, err = mediatr.Send[
			copyJawaban.CopyTemplateJawabanCommand,
			string,
		](
			ctx,
			cmdJawaban,
		)

		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	commit = true

	return result.Value.UUID.String(), nil
}
