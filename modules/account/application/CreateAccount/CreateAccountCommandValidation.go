package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateAccountCommandValidation(cmd CreateAccountCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(
			&cmd.Username,
			validation.Required.Error("Username cannot be blank"),
		),

		validation.Field(
			&cmd.Password,
			validation.Required.Error("Password cannot be blank"),
		),

		validation.Field(
			&cmd.Level,
			validation.Required.Error("Level cannot be blank"),
		),

		validation.Field(
			&cmd.Name,
			validation.Required.Error("Name cannot be blank"),
		),

		validation.Field(
			&cmd.Email,
			validation.NilOrNotEmpty.Error("Email cannot be blank"),
		),

		validation.Field(
			&cmd.Fakultas,
			validation.NilOrNotEmpty.Error("Fakultas cannot be blank"),
		),

		validation.Field(
			&cmd.Prodi,
			validation.NilOrNotEmpty.Error("Prodi cannot be blank"),
		),
	)
}
