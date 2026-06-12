package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduleTimeBankSoalCommandValidation(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	cmd := ScheduleTimeBankSoalCommand{
		UuidBankSoal: validUUID,
		TanggalMulai: "2024-01-01",
		TanggalAkhir: "2024-01-10",
		SID:          "sid-123",
		Resource:     "lpm",
	}

	// Success case
	assert.NoError(t, ScheduleTimeBankSoalCommandValidation(cmd))

	// Fail on empty UuidBankSoal
	cmdFailUuid := cmd
	cmdFailUuid.UuidBankSoal = ""
	assert.Error(t, ScheduleTimeBankSoalCommandValidation(cmdFailUuid))
}
