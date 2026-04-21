package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func DeleteTimeExtBankSoalCommandValidation(cmd DeleteTimeExtBankSoalCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.UuidBankSoal,
			validation.Required.Error("BankSoal cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
	)
}
