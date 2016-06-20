package services

type User struct {
	Id   string
	Name string
}

type UsersService interface {
	GetUser(id int) (User, error)
	UpdateUser(User) error
	DeleteUser(User) error
	CreateUser() (User, error)
}
