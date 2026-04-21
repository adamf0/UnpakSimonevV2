package application

import (
	"context"

	"UnpakSiamida/common/helper"
	domainbanksoal "UnpakSiamida/modules/banksoal/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type ScheduleTimeBankSoalCommandHandler struct {
	Repo domainbanksoal.IBankSoalRepository
}

func (h *ScheduleTimeBankSoalCommandHandler) Handle(
	ctx context.Context,
	cmd ScheduleTimeBankSoalCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidBankSoal, err := uuid.Parse(cmd.UuidBankSoal)
	if err != nil {
		return "", domainbanksoal.InvalidUuid()
	}

	var (
		existingBankSoal *domainbanksoal.BankSoalDefault
		prev             *domainbanksoal.BankSoal
	)
	g, ctxg := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		existingBankSoal, err = h.Repo.GetDefaultByUuid(ctxg, uuidBankSoal)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return domainbanksoal.NotFound(cmd.UuidBankSoal)
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		prev, err = h.Repo.GetByUuid(ctxg, uuidBankSoal)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return domainbanksoal.NotFound(cmd.UuidBankSoal)
			}
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	//[pr] validasi hanya local, simak, simpeg
	if helper.NullableString(existingBankSoal.CreatedByRef) == cmd.SID && helper.NullableString(existingBankSoal.CreatedBy) == cmd.Resource {
		result := domainbanksoal.UpdateTimeBankSoal(
			prev,
			uuidBankSoal,
			helper.StrPtr(cmd.TanggalMulai),
			helper.StrPtr(cmd.TanggalAkhir),
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		updateBankSoal := result.Value
		if err := h.Repo.Update(ctx, updateBankSoal); err != nil {
			return "", err
		}

		return result.Value.UUID.String(), nil
	} else {
		result := domainbanksoal.AddTimeBankSoalExt(
			prev,
			uuidBankSoal,
			helper.StrPtr(cmd.TanggalMulai),
			helper.StrPtr(cmd.TanggalAkhir),
			cmd.Resource,
			cmd.SID,
		)

		if !result.IsSuccess {
			return "", result.Error
		}

		createTimeBankSoal := result.Value
		if err := h.Repo.CreateExt(ctx, createTimeBankSoal); err != nil {
			return "", err
		}

		return result.Value.UUID.String(), nil
	}
}
