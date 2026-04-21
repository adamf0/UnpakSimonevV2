package application

import (
	domainKuesioner "UnpakSiamida/modules/kuesioner/domain"
	"context"
	"time"
)

type GetKuesionerJawabanByUuidQueryQueryHandler struct {
	Repo domainKuesioner.IKuesionerJawabanRepository
}

func (h *GetKuesionerJawabanByUuidQueryQueryHandler) Handle(
	ctx context.Context,
	q GetKuesionerJawabanByUuidQuery,
) ([]domainKuesioner.KuesionerJawabanDefault, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	Kuesioners, err := h.Repo.GetAllByKuesioner(
		ctx,
		q.Uuid,
	)
	if err != nil {
		return []domainKuesioner.KuesionerJawabanDefault{}, err
	}

	return Kuesioners, err
}
