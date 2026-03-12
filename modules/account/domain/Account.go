package domain

type Account struct {
	ID          string  `json:"-"`
	UUID        *string `json:"UUID"`
	Username    *string `json:"Username"`
	Password    *string `json:"-"`
	Level       *string `json:"Level"`
	Name        *string `json:"Name"`
	Email       *string `json:"Email"`
	RefFakultas *string `json:"RefFakultas"`
	Fakultas    *string `json:"Fakultas"`
	RefProdi    *string `json:"RefProdi"`
	Prodi       *string `json:"Prodi"`
	Unit        *string `json:"Unit"`
	Resource    *string `json:"Resource"`
	CodeCtx     *string `json:"CodeCtx"` //untuk membedakan dosen & mahasiswa
}

// func (Account) TableName() string {
// 	return "users"
// }
