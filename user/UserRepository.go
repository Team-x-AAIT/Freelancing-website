package user

import (
	"database/sql"
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// PsqlUserRepository is a struct that define the PsqlUserRepository type.
type PsqlUserRepository struct {
	connection *sql.DB
}

// Repository is an interface that specifies database operations on User type.
type Repository interface {
	AddUser(*entities.User) error
}

// NewPsqlUserRepository is a function that return new PsqlUserRepository type.
func NewPsqlUserRepository(conn *sql.DB) *PsqlUserRepository {
	return &PsqlUserRepository{connection: conn}
}

// AddUser is a method that adds a user to the provided database.
func (psql *PsqlUserRepository) AddUser(user *entities.User) error {
	var totalNumOfUsers int
	row := psql.connection.QueryRow("SELECT COUNT(*) FROM Users")
	row.Scan(&totalNumOfUsers)
	user.UID = fmt.Sprintf("UID%d", totalNumOfUsers+1)

	stmt, _ := psql.connection.Prepare(`INSERT INTO Users (uid, first_name, last_name, password, phonenumber, email, job_title, country, city, gender, rating)
	 																VALUES (?,?,?,?,?,?,?,?,?,?,?)`)
	_, err := stmt.Exec(
		user.UID,
		user.Firstname,
		user.Lastname,
		user.Password,
		user.Phonenumber,
		user.Email,
		user.JobTitle,
		user.Country,
		user.City,
		user.Gender,
		user.Rating)

	if err != nil {
		return err
	}
	return nil
}
