package domaintest

import (
	"UnpakSiamida/modules/prodi/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProdiDefault(t *testing.T) {
	pd := domain.ProdiDefault{
		KodeFakultas: "FKIP",
		KodeProdi:    "IPA",
		NamaProdi:    "Pendidikan IPA",
	}

	assert.Equal(t, "FKIP", pd.KodeFakultas)
	assert.Equal(t, "IPA", pd.KodeProdi)
	assert.Equal(t, "Pendidikan IPA", pd.NamaProdi)
}
