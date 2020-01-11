package service

import (
	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/Team-x-AAIT/Freelancing-website/api/myjob"
)

// struct to implement jon repo
type MyJobService struct {
	myjobRepo myjob.MyJobRepository
}

// init my job ser
func NewMyJobService(repo myjob.MyJobRepository) *MyJobService {
	return &MyJobService{myjobRepo: repo}
}

// store job
func (mjs *MyJobService) StoreMyJob(myjob *entity.MyJob) (*entity.MyJob, []error) {
	jobs, errs := mjs.myjobRepo.StoreMyJob(myjob)
	if len(errs) > 0 {
		return nil, errs
	}

	return jobs, nil
}

// gets job by id
func (mjs *MyJobService) GetMyJob(id int) ([]entity.MyJob, []error) {
	job, errs := mjs.myjobRepo.GetMyJob(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return job, nil
}
