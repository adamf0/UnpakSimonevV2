package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateTemplatePertanyaanCommandValidation(cmd UpdateTemplatePertanyaanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.UuidBankSoal,
			validation.Required.Error("BankSoal cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Pertanyaan,
			validation.Required.Error("Pertanyaan cannot be blank"),
			validation.Length(1, 5000),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.JenisPilihan,
			validation.Required.Error("Jenis pilihan cannot be blank"),
			// validation.In("radio", "checkbox", "text").Error("Invalid jenis pilihan"),
		),

		validation.Field(&cmd.Bobot,
			validation.Required.Error("Bobot cannot be blank"),
			// validation.Min(1).Error("Bobot minimal 1"),
		),

		validation.Field(&cmd.UuidKategori,
			validation.Required.Error("Kategori cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),

		validation.Field(&cmd.Required,
			validation.In(0, 1).Error("Required must be 0 or 1"),
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
