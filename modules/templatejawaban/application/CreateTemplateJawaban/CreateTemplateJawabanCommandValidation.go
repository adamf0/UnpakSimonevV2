package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateTemplateJawabanCommandValidation(cmd CreateTemplateJawabanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.UuidTemplatePertanyaan,
			validation.Required.Error("Template Pertanyaan cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),

		validation.Field(&cmd.Jawaban,
			validation.Required.Error("Jawaban cannot be blank"),
		),

		// validation.Field(&cmd.Nilai,
		// 	validation.When(cmd.IsFreeText == "0",
		// 		validation.Required.Error("Nilai wajib diisi untuk pilihan"),
		// 	).Else(
		// 		validation.Nil,
		// 	),
		// ),

		validation.Field(&cmd.IsFreeText,
			validation.Required.Error("isFreeText cannot be blank"),
			validation.In("0", "1").Error("Invalid isFreeText"),
		),

		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
		),

		validation.Field(&cmd.Resource,
			validation.Required.Error("Resource cannot be blank"),
		),
	)
}
