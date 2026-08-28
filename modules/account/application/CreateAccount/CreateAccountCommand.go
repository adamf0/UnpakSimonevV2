package application

type CreateAccountCommand struct {
	Username   string
	Password   string
	Level      string
	Name       string
	Email      *string
	EmployeeID *string
	Fakultas   *string
	Prodi      *string
}
