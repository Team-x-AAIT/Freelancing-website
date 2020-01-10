package user

import (
	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// Service is a struct that defines the UserService type.
type Service struct {
	conn IRepository
}

// Identification is a struct that hold a basic information about third party authentication.
type Identification struct {
	TpFlag          bool
	From            string
	ConfirmPassword string
}

// TempVerificationPack is a struct that contain a user type pointer and a verification token.
type TempVerificationPack struct {
	user  *entities.User
	token string
}

// ProjectUserContainer is a struct that containe both the project and owner information.
type ProjectUserContainer struct {
	Firstname     string
	Lastname      string
	Phonenumber   string
	Email         string
	JobTitle      string
	Country       string
	City          string
	Gender        string
	ProfilePic    string
	Project       *entities.Project
	CreatedString string
}

// IService is an interface that specifies what a User type can do.
type IService interface {
	RegisterUser(*entities.User) error
}
