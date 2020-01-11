package myjob

import "github.com/Team-x-AAIT/Freelancing-website/api/entity"

// my job interface
type MyJobRepository interface {
	StoreMyJob(myjob *entity.MyJob) (*entity.MyJob, []error)
	GetMyJob(id int) ([]entity.MyJob, []error)
	// UpdateJob(id int) (*entity., []error)
}