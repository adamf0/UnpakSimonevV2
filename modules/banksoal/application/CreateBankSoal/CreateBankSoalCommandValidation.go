package application

import (
	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateBankSoalCommandValidation(cmd CreateBankSoalCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Judul,
			validation.Required.Error("Judul cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
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
		validation.Field(&cmd.TanggalMulai,
			validation.Required.Error("Tanggal Mulai cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.TanggalAkhir,
			validation.Required.Error("Tanggal Akhir cannot be blank"),
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
