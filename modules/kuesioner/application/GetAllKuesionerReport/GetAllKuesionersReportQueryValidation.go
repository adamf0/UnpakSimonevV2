package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func GetAllKuesionersReportQueryValidation(q GetAllKuesionersReportQuery) error {
	return validation.ValidateStruct(&q,

		// =========================
		// JUDUL BANK SOAL
		// =========================
		validation.Field(&q.JudulBankSoal,
			validation.When(!q.Is4Year,
				validation.Required.Error("Judul bank soal wajib diisi"),
				validation.By(func(value interface{}) error {
					v, ok := value.(*string)
					if !ok || v == nil {
						return nil
					}
					if *v == "" {
						return validation.NewError("required", "Judul bank soal tidak boleh kosong")
					}
					return nil
				}),
			),
		),

		// =========================
		// SEMESTER (OPSIONAL)
		// =========================
		validation.Field(&q.Semester,
			validation.By(func(value interface{}) error {
				v, ok := value.(*string)
				if !ok || v == nil || *v == "" {
					return nil
				}

				if len(*v) > 6 {
					return validation.NewError("length", "Semester tidak valid")
				}

				return nil
			}),
		),
	)
}
