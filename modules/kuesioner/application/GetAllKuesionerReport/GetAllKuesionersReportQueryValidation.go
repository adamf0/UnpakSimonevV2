package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func GetAllKuesionersReportQueryValidation(q GetAllKuesionersReportQuery) error {
	return validation.ValidateStruct(&q,

		validation.Field(&q.JudulBankSoal,
			validation.When(q.Is4Year == false,
				validation.Required.Error("Judul bank soal wajib diisi"),
			),
			validation.When(q.JudulBankSoal != nil,
				validation.By(func(value interface{}) error {
					v, _ := value.(*string)
					if v != nil && *v == "" {
						return validation.NewError("required", "Judul bank soal tidak boleh kosong")
					}
					return nil
				}),
			),
		),

		validation.Field(&q.Semester,
			validation.By(func(value interface{}) error {
				v, _ := value.(*string)

				if v == nil || *v == "" {
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
