package repository

import (
	"database/sql"
	"errors"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// UserMockRepository is a struct that define the Repository type.
type UserMockRepository struct {
	connection *sql.DB
}

// NewUserMockRepository is a function that return new IRepository type.
func NewUserMockRepository(conn *sql.DB) *UserMockRepository {
	return &UserMockRepository{connection: conn}
}

// AddUser is a method that adds a user to the provided database.
func (psql *UserMockRepository) AddUser(user *entity.User) error {
	entity.UserMock = *user
	return nil
}

// SearchUser is a method that search for a User type from the 'users' table in the database using the provided Email address.
func (psql *UserMockRepository) SearchUser(identifier string) *entity.User {
	var userHolder entity.User
	if identifier == "UID1" || identifier == "binysimayehu@gmail.com" {
		userHolder := entity.UserMock
		return &userHolder
	}
	return &userHolder
}

// UpdateUser is a method that updates the user data in the database using the provided parameters.
func (psql *UserMockRepository) UpdateUser(user *entity.User) error {

	if user.UID == "UID1" {
		return nil
	}
	return errors.New("user not found")

}

// RegisterTPUsers is a method that adds a third party user to the provided database.
func (psql *UserMockRepository) RegisterTPUsers(uid, email, password, from string) error {
	return nil
}

// UpdateTPUsers is a method that updates the third party user data in the database using the provided parameters.
func (psql *UserMockRepository) UpdateTPUsers(uid, email string) error {
	return nil
}

// SearchTPUser is a method that search for a User type from the 'tp_users' table in database using the provided Email address.
func (psql *UserMockRepository) SearchTPUser(identifier string) (*entity.User, string) {

	var userHolder entity.User
	if identifier == "UID1" || identifier == "binysimayehu@gmail.com" {
		userHolder := entity.UserMock
		return &userHolder, "UID1"
	}
	return &userHolder, ""

}

// RemoveUser remove a user from user table.
func (psql *UserMockRepository) RemoveUser(uid string) (*entity.User, error) {

	var userHolder entity.User
	if uid == "UID1" || uid == "binysimayehu@gmail.com" {
		userHolder := entity.UserMock
		return &userHolder, nil
	}
	return &userHolder, errors.New("user not found")
}

// RemoveTPUser remove a user from user_tp table.
func (psql *UserMockRepository) RemoveTPUser(uid string) (*entity.User, error) {

	var userHolder entity.User
	if uid == "UID1" || uid == "binysimayehu@gmail.com" {
		userHolder := entity.UserMock
		return &userHolder, nil
	}
	return &userHolder, errors.New("user not found")
}

// CountMember is a method that is used for counting the member of a table where our table name is provided as an argument.
func (psql *UserMockRepository) CountMember(tableName string) (totalNumOfMembers int) {

	return 1
}

// AddMatchTag is a method that adds a matchTag to the match tag table.
func (psql *UserMockRepository) AddMatchTag(uid string, category string, subcategory string, worktype int) error {

	entity.MatchTagMock.UID = uid
	entity.MatchTagMock.Category = category
	entity.MatchTagMock.Subcategory = subcategory
	entity.MatchTagMock.WorkType = worktype

	return nil
}

// RemoveMatchTag is a method that removes a matchTag from the match tag table.
func (psql *UserMockRepository) RemoveMatchTag(uid string, category string, subcategory string, worktype int) error {

	if entity.MatchTagMock.UID == uid && entity.MatchTagMock.Category ==
		category && entity.MatchTagMock.Subcategory == subcategory && worktype == entity.MatchTagMock.WorkType {
		return nil
	}

	return errors.New("matchtag not found")
}

// GetUserMatchTags is a method that returns all the match tag the user defined.
func (psql *UserMockRepository) GetUserMatchTags(uid string) []*entity.MatchTag {

	if entity.MatchTagMock.UID == uid {
		return []*entity.MatchTag{&entity.MatchTagMock}
	}
	return []*entity.MatchTag{&entity.MatchTag{}}
}

// SearchProjectWMatchTag is a method that returns all the projects that match the provided tag.
func (psql *UserMockRepository) SearchProjectWMatchTag(matchTag *entity.MatchTag) []*entity.Project {
	return []*entity.Project{&entity.Project{}}
}

// RemoveFileDB is a method that is used for clearing a certain column.
func (psql *UserMockRepository) RemoveFileDB(uid string, columnName string) string {
	return ""
}

// RemoveFile is a method that removes a given file path from the assets folder.
func (psql *UserMockRepository) RemoveFile(filename string, folder string) error {

	return nil
}

// SearchMember is a method that is used for searching the member of a table where our table name is provided as an argument.
func (psql *UserMockRepository) SearchMember(tableName string, columnValue string) bool {

	return true
}

// GetOwner is a method that returns the owner of a project by searching through the project user table.
func (psql *UserMockRepository) GetOwner(pid string) string {

	if pid == "PID1" {
		return entity.UserMock.UID
	}
	return ""
}
