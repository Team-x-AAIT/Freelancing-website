package application

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// IService is an interface that specifies what a application type can do.
type IService interface {
	Apply(string, string, string) error
	GetApplication(applicantUID, pid string) (*entity.ApplicationBag, error)
	GetProjectApplicantsID(string) []*entity.ApplicationBag
	GetUserApplicationHistory(string) []*entity.ApplicationBag
	HireApplicant(string, string) error
	RemoveApplicant(string, string) (*entity.ApplicationBag, error)
	RemoveApplicationInfo(applicantUID string, pid string) error
	CheckApplicationStatus(*entity.Project, string) int64
	// GetAppliedFor(uid string) []*entity.ApplicationBag
	// GetHiredFor(uid string) []*entity.ApplicationBag
}
