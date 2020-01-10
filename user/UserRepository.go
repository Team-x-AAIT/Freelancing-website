package user

import (
	"database/sql"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// Repository is a struct that define the Repository type.
type Repository struct {
	connection *sql.DB
}

// MatchTag is a struct that define a type of a match tag returned from database.
type MatchTag struct {
	UID         string
	Category    string
	Subcategory string
	WorkType    int
}

// IRepository is an interface that specifies database operations on User type.
type IRepository interface {
	AddUser(*entities.User) error
}
