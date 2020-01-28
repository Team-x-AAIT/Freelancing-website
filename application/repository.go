package application

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// IRepository is an interface that specifies database operations on project type.
type IRepository interface {
	AddApplication(string, string, string) error
	AddApplicationToHistory(string, string, string) error
	GetApplicants(string) []*entity.ApplicationBag
	GetApplication(string, string) (*entity.ApplicationBag, error)
	GetApplicationFromHistory(string, string) *entity.ApplicationBag
	GetUserApplicationHistory(string) []*entity.ApplicationBag
	HireApplicant(string, string) error
	UpdateApplicationTable(string, string, string, bool) error
	UpdateApplicationHistoryTable(string, string, string, bool, bool) error
	RemoveApplicationInfo(string, string) error
	RemoveUnHiredApplicants(string) error
}
