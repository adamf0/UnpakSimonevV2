package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ScheduleTimeBankSoalCommandValidation(cmd ScheduleTimeBankSoalCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.UuidBankSoal,
			validation.Required.Error("BankSoal cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		// validation.Field(&cmd.Content,
		// 	validation.Required.Error("Content cannot be blank"),
		// ),
		// validation.Field(&cmd.Deskripsi,
		// 	validation.Required.Error("Deskripsi cannot be blank"),
		// ),
		validation.Field(&cmd.TanggalMulai,
			validation.Required.Error("Tanggal Mulai cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalAkhir,
			validation.Required.Error("Tanggal Akhir cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Resource,
			validation.Required.Error("Resource cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}
