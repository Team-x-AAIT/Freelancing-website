package repository

import (
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// conn struct
type MyJobGormRepo struct {
	conn *gorm.DB
}

// init conn struct
func NewMyJobGormRepo(dbconn *gorm.DB) *MyJobGormRepo {
	return &MyJobGormRepo{conn: dbconn}
}

// store job repo layer
func (mr *MyJobGormRepo) StoreMyJob(myjob *entity.MyJob) (*entity.MyJob, []error) {
	job := myjob
	errs := mr.conn.Create(job).GetErrors()

	for _, err := range errs {
		pqerr := err.(*pq.Error)
		fmt.Println(pqerr)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return job, nil
}

// gets job by id
func (mr *MyJobGormRepo) GetMyJob(id int) ([]entity.MyJob, []error) {
	jobs := []entity.MyJob{}
	errs := mr.conn.Where("user_id = ?", id).Find(&jobs).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return jobs, nil
}
