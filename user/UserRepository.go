package user

import (
	"database/sql"

	"github.com/Team-x-AAIT/Freelancing-website/entites"
)

// PsqlUserRepository is a struct that define the PsqlUserRepository type.
type PsqlUserRepository struct {
	connection *sql.Conn
}

// UserRepository is an interface that specifies database operations on User type.
type UserRepository interface {
	AddUser(entites.User) error
}

// NewPsqlUserRepository is a function that return new PsqlUserRepository type.
func NewPsqlUserRepository(conn *sql.Conn) *PsqlUserRepository {
	return &PsqlUserRepository{connection: conn}
}

// AddUser is a method that adds a user to the provided database.
func (pr *PsqlUserRepository) AddUser(user entites.User) (err error) {
	err = nil
	return
}
