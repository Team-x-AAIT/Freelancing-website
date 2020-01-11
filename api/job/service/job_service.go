package service

import (
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/Team-x-AAIT/Freelancing-website/api/job"
)

// ser
type JobService struct {
	jobRepo job.JobRepository
}

// init ser
func NewJobService(repo job.JobRepository) *JobService {
	return &JobService{jobRepo: repo}
}

// calls jobs on repository layer
func (js *JobService) Jobs(category string) ([]entity.Job, []error) {
	jobs, errs := js.jobRepo.Jobs(category)
	if len(errs) > 0 {
		return nil, errs
	}

	return jobs, nil
}

// job
func (js *JobService) Job(job *entity.Job) (*entity.Job, []error) {
	job, errs := js.jobRepo.Job(job)
	fmt.Println(errs)
	if len(errs) > 0 {
		return nil, errs
	}

	return job, nil
}

// it returns ajob after it feated it by its id
func (js *JobService) JobByID(id uint) (*entity.Job, []error) {
	job, errs := js.jobRepo.JobByID(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return job, nil
}

// update
func (js *JobService) UpdateJob(job *entity.Job) (*entity.Job, []error) {
	myjob, errs := js.jobRepo.UpdateJob(job)

	if len(errs) > 0 {
		return nil, errs
	}

	return myjob, nil
}

// delets job calls the delete job method on repository layer
func (js *JobService) DeleteJob(id uint) (*entity.Job, []error) {
	usr, errs := js.jobRepo.DeleteJob(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// stores user
func (js *JobService) StoreJob(job *entity.Job) (*entity.Job, []error) {
	myjob, errs := js.jobRepo.StoreJob(job)
	if len(errs) > 0 {
		return nil, errs
	}

	return myjob, nil
}
