package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SaveKuesionerJawabanCommandValidation(cmd SaveKuesionerJawabanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.UuidKuesioner,
			validation.Required.Error("Kuesioner cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.UuidPertanyaan,
			validation.Required.Error("Pertanyaan cannot be blank"),
		),
		validation.Field(&cmd.Jawaban,
			validation.Required.Error("Jawaban cannot be blank"),
			validation.By(helper.ValidateJSONArray),
		),
		validation.Field(&cmd.SID,
			validation.Required.Error("SID cannot be blank"),
		),
		validation.Field(&cmd.Resource,
			validation.Required.Error("Resource cannot be blank"),
		),
	)
}
