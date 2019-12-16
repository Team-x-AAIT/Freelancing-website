package user

import (
	"database/sql"
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// Repository is a struct that define the Repository type.
type Repository struct {
	connection *sql.DB
}

// IRepository is an interface that specifies database operations on User type.
type IRepository interface {
	AddUser(*entities.User) error
	SearchUser(string) *entities.User
	RegisterTPUsers(string, string, string, string) error
	SearchTPUser(string) (*entities.User, string)
	CountMember(string) int
}

// NewRepository is a function that return new IRepository type.
func NewRepository(conn *sql.DB) IRepository {
	return &Repository{connection: conn}
}

// AddUser is a method that adds a user to the provided database.
func (psql *Repository) AddUser(user *entities.User) error {

	totalNumOfUsers := psql.CountMember("users")

	user.UID = fmt.Sprintf("UID%d", totalNumOfUsers+1)

	stmt, _ := psql.connection.Prepare(`INSERT INTO Users (uid,first_name, last_name, password, phonenumber, email, 
		job_title, country, city, gender,cv, profile_pic, bio, rating) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
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
		user.CV,
		user.ProfilePic,
		user.Bio,
		user.Rating)

	if err != nil {
		return err
	}
	return nil
}

// SearchUser is a method that search for a User type from the 'users' table in the database using the provided Email address.
func (psql *Repository) SearchUser(identifier string) *entities.User {
	stmt, err := psql.connection.Prepare("SELECT * FROM Users WHERE email=? || uid=?")
	row := stmt.QueryRow(identifier, identifier)
	if err != nil {
		panic(err)
	}
	var user entities.User
	row.Scan(
		&user.UID,
		&user.Firstname,
		&user.Lastname,
		&user.Password,
		&user.Phonenumber,
		&user.Email,
		&user.JobTitle,
		&user.Country,
		&user.City,
		&user.Gender,
		&user.CV,
		&user.ProfilePic,
		&user.Bio,
		&user.Rating)
	return &user
}

// RegisterTPUsers is a method that adds a third party user to the provided database.
func (psql *Repository) RegisterTPUsers(uid, email, password, from string) error {
	stmt, _ := psql.connection.Prepare(`INSERT INTO tp_users (uid, email, password, origin) VALUES (?,?,?,?)`)
	_, err := stmt.Exec(uid, email, password, from)

	if err != nil {
		return err
	}
	return nil
}

// SearchTPUser is a method that search for a User type from the 'tp_users' table in database using the provided Email address.
func (psql *Repository) SearchTPUser(identifier string) (*entities.User, string) {

	stmt, err := psql.connection.Prepare("SELECT * FROM tp_users WHERE email=? || uid=?")
	row := stmt.QueryRow(identifier, identifier)
	if err != nil {
		panic(err)
	}
	var from string
	var user entities.User
	row.Scan(
		&user.UID,
		&user.Email,
		&user.Password,
		&from)
	return &user, from
}

// CountMember is a method that is used for counting the member of a table where our table name is provided as an argument.
func (psql *Repository) CountMember(tableName string) (totalNumOfMembers int) {

	stmt, err := psql.connection.Prepare("SELECT COUNT(*) FROM " + tableName)
	if err != nil {
		return
	}
	row := stmt.QueryRow()
	row.Scan(&totalNumOfMembers)
	return

}
