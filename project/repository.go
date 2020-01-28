package project

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// IRepository is an interface that specifies database operations on project type.
type IRepository interface {
	GetProject(string) *entity.Project
	AddProject(*entity.Project) (string, error)
	UpdateProject(*entity.Project) (string, error)
	RemoveProject(string) error
	SearchProject(string, string, int64, float64, float64, int64) []*entity.Project
	MarkAsClosed(string) error

	AttachFiles(string, string) error
	GetAttachedFiles(string) []string
	RemoveAttachedFiles(string) error
	RemoveAttachedFile(string, string) error
	RemoveFile(string) error
	CountMember(string) int
	SearchMember(string, string) bool

	GetLinkedProjects(string) []string
	LinkProject(string, string) error
	UnLinkProject(string, string) error
	SearchLink(string, string) bool

	GetCategories() []string
	GetSubCategories() []string
	GetSubCategoriesOf(category string) []string
}
