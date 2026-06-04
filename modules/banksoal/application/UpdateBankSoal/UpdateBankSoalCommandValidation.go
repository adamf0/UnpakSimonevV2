package application

import (
	helper "UnpakSiamida/common/helper"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateBankSoalCommandValidation(cmd UpdateBankSoalCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
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
			validation.Match(regexp.MustCompile(`^\d{4}(01|02)$`)).
				Error("Semester invalid format"),
		),
		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
		),
		validation.Field(&cmd.Resource,
			validation.Required.Error("Resource cannot be blank"),
		),
	)
}
