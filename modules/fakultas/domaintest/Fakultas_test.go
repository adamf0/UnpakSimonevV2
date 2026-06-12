package domaintest

import (
	"UnpakSiamida/modules/fakultas/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakultasDefault(t *testing.T) {
	fd := domain.FakultasDefault{
		KodeFakultas: "FKIP",
		NamaFakultas: "Keguruan dan Ilmu Pendidikan",
	}

	assert.Equal(t, "FKIP", fd.KodeFakultas)
	assert.Equal(t, "Keguruan dan Ilmu Pendidikan", fd.NamaFakultas)
}
