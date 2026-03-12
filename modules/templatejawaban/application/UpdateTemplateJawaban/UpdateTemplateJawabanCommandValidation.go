package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateTemplateJawabanCommandValidation(cmd UpdateTemplateJawabanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.UuidTemplatePertanyaan,
			validation.Required.Error("Template Pertanyaan cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Jawaban,
			validation.Required.Error("Jawaban cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Nilai,
			validation.Required.Error("Nilai cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.IsFreeText,
			validation.Required.Error("isFreeText cannot be blank"),
			validation.In("0", "1").Error("Invalid isFreeText"),
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
