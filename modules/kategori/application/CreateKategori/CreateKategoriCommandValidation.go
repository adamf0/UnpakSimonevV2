package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateKategoriCommandValidation(cmd CreateKategoriCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.NamaKategori,
			validation.Required.Error("NamaKategori cannot be blank"),
		),
		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
		),
		validation.Field(&cmd.Resource,
			validation.Required.Error("Resource cannot be blank"),
		),
	)
}
