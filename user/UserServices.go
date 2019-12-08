package user

import "github.com/Team-x-AAIT/Freelancing-website/entites"

// UserService is a struct that defines the UserService type.
type UserService struct {
	conn *UserRepository
}

// UserServices is an interface that specifies what User type can do.
type UserServices interface {
	RegisterUser(entites.User) error
}

// NewUserService is a function that returns a new UserService type.
func NewUserService(connection *UserRepository) *UserService {
	return &UserService{conn: connection}
}

// RegisterUser is a method that register a new user to the system.
func (service *UserService) RegisterUser(user entites.User) (err error) {
	err = nil
	return
}
