package user

import "github.com/Team-x-AAIT/Freelancing-website/entities"

// Service is a struct that defines the UserService type.
type Service struct {
	conn Repository
}

// Services is an interface that specifies what User type can do.
type Services interface {
	RegisterUser(*entities.User) error
}

// NewUserService is a function that returns a new UserService type.
func NewUserService(connection Repository) Services {
	return &Service{conn: connection}
}

// RegisterUser is a method that register a new user to the system.
func (service *Service) RegisterUser(user *entities.User) (err error) {
	err = service.conn.AddUser(user)

	if err != nil {
		return err
	}
	return nil
}
