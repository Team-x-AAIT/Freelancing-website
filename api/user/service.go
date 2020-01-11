package user

import "github.com/Team-x-AAIT/Freelancing-website/api/entity"

// user service interface
type UserService interface {
	Users() ([]entity.User, []error)
	User(user *entity.User) (*entity.User, []error)
	UserByID(id uint) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
	RecommendedJobs(id uint) ([]entity.Job, []error)
}
