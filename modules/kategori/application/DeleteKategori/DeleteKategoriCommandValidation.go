package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func DeleteKategoriCommandValidation(cmd DeleteKategoriCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Uuid,
			validation.Required.Error("mode cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	)
}
