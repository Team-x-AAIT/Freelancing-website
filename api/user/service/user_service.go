package service

import (
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/Team-x-AAIT/Freelancing-website/api/user"
)
// user repo iplemnted
type UserService struct {
	userRepo user.UserRepository
}
// init user service
func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}
// calls users in repo
func (us *UserService) Users() ([]entity.User, []error) {
	users, errs := us.userRepo.Users()
	if len(errs) > 0 {
		return nil, errs
	}

	return users, nil
}
// calls user on repo layer
func (us *UserService) User(user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.User(user)
	fmt.Println(errs)
	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// UserByID returns a single user by its id
func (us *UserService) UserByID(id uint) (*entity.User, []error) {
	usr, errs := us.userRepo.UserByID(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// UpdateUser updates user
func (us *UserService) UpdateUser(user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.UpdateUser(user)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// DeleteUser deletes a single  user
func (us *UserService) DeleteUser(id uint) (*entity.User, []error) {
	usr, errs := us.userRepo.DeleteUser(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// StoreUser will insert a new  user
func (us *UserService) StoreUser(user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.StoreUser(user)
	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

func (us *UserService) RecommendedJobs(id uint) ([]entity.Job, []error) {

	job, errs := us.userRepo.RecommendedJobs(id)
	if len(errs) > 0 {
		return nil, errs
	}

	return job, nil
}
