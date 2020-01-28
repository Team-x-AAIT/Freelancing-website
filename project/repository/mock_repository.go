package repository

import (
	"database/sql"
	"errors"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// MockRepository is a struct that define the Repository type.
type MockRepository struct {
	connection *sql.DB
}

// NewMockRepository is a function that return new Repository type.
func NewMockRepository(conn *sql.DB) *MockRepository {
	return &MockRepository{connection: conn}
}

// GetProject is a method that retrive the fake project object.
func (mock *MockRepository) GetProject(pid string) *entity.Project {

	var project entity.Project
	if pid == "PID1" {
		project := entity.ProjectMock
		return &project
	}
	return &project

}

// AddProject is a method that adds a project to the provided database.
func (mock *MockRepository) AddProject(project *entity.Project) (string, error) {

	project.ID = entity.ProjectMock.ID
	return project.ID, nil
}

// UpdateProject is a method that updates a project using the project id.
func (mock *MockRepository) UpdateProject(project *entity.Project) (string, error) {

	if project.ID == "PID1" {
		entity.ProjectMock = *project
		return project.ID, nil
	}
	return "", errors.New("error has occured")
}

// RemoveProject is a method that remove a project from the database using the project id.
func (mock *MockRepository) RemoveProject(pid string) error {

	if pid == "PID1" {
		return nil
	}
	return errors.New("not found")
}

// SearchProject is a method that is used for searching a projects from the database using a search-key.
func (mock *MockRepository) SearchProject(searchKey string, searchBy string, filterType int64, filterValue1 float64, filterValue2 float64, pageNumber int64) []*entity.Project {

	projects := []*entity.Project{&entity.ProjectMock}
	return projects
}

// MarkAsClosed is a method that mark a project as closed by updating the closed filed of a project in the database.
func (mock *MockRepository) MarkAsClosed(pid string) error {

	if pid == "PID1" {
		entity.ProjectMock.Closed = true
		return nil
	}

	return errors.New("not found")
}

// AttachFiles is a method that adds a filename and its owner that is the project id into attached_files table.
func (mock *MockRepository) AttachFiles(pid string, fileName string) error {

	if pid == "PID1" {
		return nil
	}

	return errors.New("not found")

}

// GetAttachedFiles is a method that returns names of the files that are attached to the project.
func (mock *MockRepository) GetAttachedFiles(pid string) []string {

	if pid == "PID1" {
		return entity.ProjectMock.AttachedFiles
	}

	return []string{}
}

// RemoveAttachedFiles is a method that is used for removing all the files attached to a project from the database.
func (mock *MockRepository) RemoveAttachedFiles(pid string) error {

	if pid == "PID1" {
		entity.ProjectMock.AttachedFiles = []string{}
		return nil
	}

	return errors.New("not found")
}

// RemoveAttachedFile is a method that is used for removing a file attached to a project from the database.
func (mock *MockRepository) RemoveAttachedFile(pid string, fileName string) error {

	attachedFiles := []string{}
	if pid == "PID1" {

		for _, value := range entity.ProjectMock.AttachedFiles {
			if value != fileName {
				attachedFiles = append(attachedFiles, value)
			}
		}

		entity.ProjectMock.AttachedFiles = attachedFiles
	}

	return errors.New("not found")
}

// RemoveFile is a method that removes a given file path from the assets folder.
func (mock *MockRepository) RemoveFile(filename string) error {

	return nil
}

// CountMember is a method that is used for counting the member of a table where our table name is provided as an argument.
func (mock *MockRepository) CountMember(tableName string) (totalNumOfMembers int) {

	if tableName == "projects" {
		totalNumOfMembers = 1
	}
	return
}

// SearchMember is a method that is used for searching the member of a table where our table name is provided as an argument.
func (mock *MockRepository) SearchMember(tableName string, columnValue string) bool {

	if tableName == "projects" && columnValue == "PID1" {
		return true
	}
	return false

}

// GetLinkedProjects is a method that returns all the linked projects id to a certain user.
func (mock *MockRepository) GetLinkedProjects(uid string) []string {

	if uid == "UID1" {
		return []string{entity.ProjectMock.ID}
	}

	return []string{}
}

// LinkProject is a method that links a project to user.
func (mock *MockRepository) LinkProject(uid, pid string) error {

	if pid == "PID1" && uid == "UID1" {
		return nil
	}
	return errors.New("link not found")
}

// UnLinkProject is a method that a unlinks a project from a user.
func (mock *MockRepository) UnLinkProject(uid, pid string) error {

	if pid == "PID1" && uid == "UID1" {
		return nil
	}
	return errors.New("link not found")
}

// SearchLink is a method that is used for identifying a certain link exists between a project and a user.
func (mock *MockRepository) SearchLink(uid, pid string) bool {

	if uid == "UID1" && pid == "PID1" {
		return true
	}
	return false
}

// GetCategories retrives all the categories stored in the database.
func (mock *MockRepository) GetCategories() []string {
	return []string{"Mechanical Systems Engineer", "Help Desk Operator", "Dental Hygienist"}
}

// GetSubCategories retrives all the sub categories stored in the database.
func (mock *MockRepository) GetSubCategories() []string {
	return []string{"Mechanical Systems Engineer", "Help Desk Operator", "Dental Hygienist"}
}

// GetCategoryID retrives the category id of a certain category.
func (mock *MockRepository) GetCategoryID(category string) int {

	if category == "Mechanical Systems Engineer" {
		return 1
	}

	return -1

}

// GetSubCategoriesOf retrives the sub categories of a given category.
func (mock *MockRepository) GetSubCategoriesOf(category string) []string {

	if category == "Mechanical Systems Engineer" {
		return []string{"Mechanical Systems Engineer"}
	}
	return []string{}

}
