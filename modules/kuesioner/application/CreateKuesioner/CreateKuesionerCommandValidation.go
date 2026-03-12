package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateKuesionerCommandValidation(cmd CreateKuesionerCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.UuidBankSoal,
			validation.Required.Error("BankSoal cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Tanggal,
			validation.Required.Error("Tanggal cannot be blank"),
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
