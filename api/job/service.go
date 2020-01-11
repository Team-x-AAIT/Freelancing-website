package job

import "github.com/Team-x-AAIT/Freelancing-website/api/entity"

// job service interface
type JobService interface {
	StoreJob(job *entity.Job) (*entity.Job, []error)
	Jobs(search string) ([]entity.Job, []error)
	Job(job *entity.Job) (*entity.Job, []error)
	JobByID(id uint) (*entity.Job, []error)
	UpdateJob(user *entity.Job) (*entity.Job, []error)
	DeleteJob(id uint) (*entity.Job, []error)
}
