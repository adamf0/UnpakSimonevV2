package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func WhoamiCommandValidation(cmd WhoamiCommand) error {
	if (cmd.SID == nil || *cmd.SID == "") &&
		(cmd.NIM == nil || *cmd.NIM == "") &&
		(cmd.NIDN == nil || *cmd.NIDN == "") &&
		(cmd.NIP == nil || *cmd.NIP == "") {
		return validation.NewError("validation_whoami", "SID, NIM, NIDN, or NIP must be provided")
	}
	return nil
}
