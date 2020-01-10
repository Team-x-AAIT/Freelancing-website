package project

import (
	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// Service is a struct that defines the Project Service type.
type Service struct {
	conn IRepository
}

// IService is an interface that specifies what a Project type can do.
type IService interface {
	PostProject(*entities.Project, string) (string, error)
	SearchProjectByID(string) *entities.Project
}
