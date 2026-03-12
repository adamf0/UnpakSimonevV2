package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("TemplateJawaban.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("TemplateJawaban.InvalidUuid", "uuid is invalid")
}

func InvalidTemplatePertanyaan() domain.Error {
	return domain.NotFoundError("TemplateJawaban.InvalidTemplatePertanyaan", "kategori is invalid")
}

func NotFoundTemplatePertanyaan() domain.Error {
	return domain.NotFoundError("TemplateJawaban.NotFoundTemplatePertanyaan", "kategori is not found")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("TemplateJawaban.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("TemplateJawaban.NotFound", fmt.Sprintf("TemplateJawaban with identifier %s not found", id))
}

func InvalidOwner() domain.Error {
	return domain.NotFoundError("TemplateJawaban.InvalidOwner", "only lpm / fakultas / prodi can create bank soal")
}
