package application

import "gorm.io/gorm"

type CopyTemplateJawabanCommand struct {
	Tx                         *gorm.DB
	SourceTemplatePertanyaanID uint
	TargetTemplatePertanyaanID uint
}
