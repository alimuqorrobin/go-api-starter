package user

type Repository interface {
	FindByID(id uint) (*User, error)
	Create(user *User) error
}
