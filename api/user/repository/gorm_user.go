package repository

import (
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// gotm conn
type UserGormRepo struct {
	conn *gorm.DB
}

// init conn
func NewUserGormRepo(dbconn *gorm.DB) *UserGormRepo {
	return &UserGormRepo{conn: dbconn}
}

// gets users
func (ur *UserGormRepo) Users() ([]entity.User, []error) {
	users := []entity.User{}
	errs := ur.conn.Find(&users).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}

	return users, nil
}

// return a single user for login purpose
func (ur *UserGormRepo) User(user *entity.User) (*entity.User, []error) {
	usr := entity.User{}
	errs := ur.conn.Where("email = ? AND password = ?", user.Email, user.Password).First(&usr).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}

	return &usr, nil
}

// find user by id
func (ur *UserGormRepo) UserByID(id uint) (*entity.User, []error) {
	usr := entity.User{}
	errs := ur.conn.First(&usr, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}

	return &usr, nil
}

// updates user
func (ur *UserGormRepo) UpdateUser(user *entity.User) (*entity.User, []error) {
	usr := user
	errs := ur.conn.Save(usr).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// delete user by id
func (ur *UserGormRepo) DeleteUser(id uint) (*entity.User, []error) {
	usr, errs := ur.UserByID(id)

	if len(errs) > 0 {
		return nil, errs
	}

	errs = ur.conn.Delete(usr, id).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// stores user
func (ur *UserGormRepo) StoreUser(user *entity.User) (*entity.User, []error) {
	usr := user
	errs := ur.conn.Create(usr).GetErrors()

	for _, err := range errs {
		pqerr := err.(*pq.Error)
		fmt.Println(pqerr)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

func (ur *UserGormRepo) RecommendedJobs(id uint) ([]entity.Job, []error) {
	job := []entity.Job{}
	// errs := ur.conn.Create(usr).GetErrors()
	errs := ur.conn.Table("jobs").Select("id,title,description,category,jobs.user_id,created_at").Joins("left join my_jobs on category = job").Find(&job).GetErrors()
	// .Where("agent_id = ? AND status = 'pending'",id).Find(&services).GetErrors()
	fmt.Println(job)
	if len(errs) > 0 {
		return nil, errs
	}
	return job, nil
}
