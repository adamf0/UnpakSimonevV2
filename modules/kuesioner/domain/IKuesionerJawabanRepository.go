package domain

import (
	"context"

	"gorm.io/gorm"
)

type IKuesionerJawabanRepository interface {
	GetByPertanyaanAndUser(
		ctx context.Context,
		pertanyaanID uint,
		sid string,
		resource string,
	) ([]KuesionerJawaban, error)
	GetTotalInputByKuesionerIDs(
		ctx context.Context,
		ids []uint,
	) (map[string]uint, error)
	GetAllByKuesioner(
		ctx context.Context,
		uuidkuesioner string,
	) ([]KuesionerJawabanDefault, error)

	Create(ctx context.Context, data *KuesionerJawaban) error
	Delete(ctx context.Context, id uint) error

	WithTx(tx any) IKuesionerJawabanRepository
	BeginTx(ctx context.Context) (*gorm.DB, error)
}
