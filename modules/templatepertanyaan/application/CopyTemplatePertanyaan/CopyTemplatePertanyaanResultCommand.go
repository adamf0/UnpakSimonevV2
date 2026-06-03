package application

import "gorm.io/gorm"

type CopyTemplatePertanyaanResultCommand struct {
	Tx               *gorm.DB
	SourceBankSoalID uint
	TargetBankSoalID uint
	Resource         string
	Sid              string
}
