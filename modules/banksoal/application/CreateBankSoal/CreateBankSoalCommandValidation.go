package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateBankSoalCommandValidation(cmd CreateBankSoalCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Judul,
			validation.Required.Error("Judul cannot be blank"),
		),
		// validation.Field(&cmd.Content,
		// 	validation.Required.Error("Content cannot be blank"),
		// ),
		// validation.Field(&cmd.Deskripsi,
		// 	validation.Required.Error("Deskripsi cannot be blank"),
		// ),
		validation.Field(&cmd.Semester,
			validation.Required.Error("Semester cannot be blank"),
		),
		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
		),
		validation.Field(&cmd.Resource,
			validation.Required.Error("Resource cannot be blank"),
		),
	)
}
