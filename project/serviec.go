package project

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// IService is an interface that specifies what a Project type can do.
type IService interface {
	ViewProject(pid string, owner *entity.User, projectW *entity.Project) *entity.ProjectUserContainer
	PostProject(*entity.Project, string) (string, entity.ErrorBag)
	SearchProjectByID(string) *entity.Project
	UpdateProject(*entity.Project) (string, entity.ErrorBag)
	RemoveProjectInformation(string, string) (*entity.Project, error)
	FindProject(string, string, int64, float64, float64, int64) []*entity.Project
	GetSentProjects(uid string) []*entity.Project

	RemoveAttachedFile(string, string) error
	RemoveFile(string) error
	SearchMember(string, string) bool
	CountMember(string) int
	AttachFiles(string, string) error
	SearchLink(uid string, pid string) bool

	GetCategories() []string
	GetSubCategories() []string
	GetSubCategoriesOf(category string) []string
}
