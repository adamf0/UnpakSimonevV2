package application

type UpdateAccountCommand struct {
	Uuid       string
	Username   string
	Password   *string
	Level      string
	Name       string
	Email      *string
	EmployeeID *string
	Fakultas   *string
	Prodi      *string
}
