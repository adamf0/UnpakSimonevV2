package application

import (
	"context"

	domainKategori "UnpakSiamida/modules/kategori/domain"
)

type SetupUuidKategoriCommandHandler struct {
	Repo domainKategori.IKategoriRepository
}

func (h *SetupUuidKategoriCommandHandler) Handle(
	ctx context.Context,
	cmd SetupUuidKategoriCommand,
) (string, error) {

	err := h.Repo.SetupUuid(ctx)
	if err != nil {
		return "", err
	}

	return "berhasil setup uuid pada data", nil
}
