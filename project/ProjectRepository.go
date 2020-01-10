package project

import (
	"database/sql"
	"time"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// Repository is a struct that define the Repository type.
type Repository struct {
	connection *sql.DB
}

// ApplicationBag is a struct that holds the application form or result
type ApplicationBag struct {
	PID         string
	ApplicantID string
	Proposal    string
	Hired       bool
	Seen        bool
	Status      int64
	CreatedAt   time.Time
	Project     *entities.Project
}

// IRepository is an interface that specifies database operations on project type.
type IRepository interface {
	AddProject(*entities.Project) (string, error)
	CountMember(string) int
	SearchMember(string, string) bool
}
