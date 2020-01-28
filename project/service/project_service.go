package service

import (
	"errors"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/project"
)

// Service is a struct that defines the Project Service type.
type Service struct {
	conn project.IRepository
}

// NewProjectService is a function that returns a new project service type.
func NewProjectService(connection project.IRepository) *Service {
	return &Service{conn: connection}
}

// PostProject is a method that validates and adds a project to the system.
func (service *Service) PostProject(project *entity.Project, uid string) (string, entity.ErrorBag) {

	if err := service.validateProject(project); len(err) > 0 {
		return "", err
	}

	pid, err := service.conn.AddProject(project)
	if err != nil {
		var errMap entity.ErrorBag = make(map[string]string)
		errMap["failed"] = "Unable to add project!"
		return "", errMap
	}

	err = service.conn.LinkProject(uid, pid)
	if err != nil {
		var errMap entity.ErrorBag = make(map[string]string)
		errMap["failed"] = "Unable to link project!"
		return "", errMap
	}

	return pid, nil
}

// SearchProjectByID is a method that returns projects using the project id.
func (service *Service) SearchProjectByID(pid string) *entity.Project {

	project := service.conn.GetProject(pid)
	attachedFiles := service.conn.GetAttachedFiles(pid)
	project.AttachedFiles = attachedFiles
	return project

}

// ViewProject is a method that returns a project with its owner information.
func (service *Service) ViewProject(pid string, owner *entity.User, projectW *entity.Project) *entity.ProjectUserContainer {

	layoutUS := "January 2, 2006"

	projectUserContainer := new(entity.ProjectUserContainer)
	projectUserContainer.Firstname = owner.Firstname
	projectUserContainer.Lastname = owner.Lastname
	projectUserContainer.Phonenumber = owner.Phonenumber
	projectUserContainer.Email = owner.Email
	projectUserContainer.JobTitle = owner.JobTitle
	projectUserContainer.Country = owner.Country
	projectUserContainer.City = owner.City
	projectUserContainer.Gender = owner.Gender
	projectUserContainer.ProfilePic = owner.ProfilePic
	projectUserContainer.Project = projectW
	projectUserContainer.CreatedString = projectUserContainer.Project.CreatedAt.Format(layoutUS)

	return projectUserContainer
}

// UpdateProject is a method that is used for updating a project profile.
func (service *Service) UpdateProject(project *entity.Project) (string, entity.ErrorBag) {

	if err := service.validateProject(project); len(err) > 0 {
		return "", err
	}

	pid, err := service.conn.UpdateProject(project)
	if err != nil {
		var errMap entity.ErrorBag = make(map[string]string)
		errMap["failed"] = "Unable to update project!"
		return "", errMap
	}
	return pid, nil
}

// RemoveProjectInformation is a method that is used for removing project and its dependencies.
func (service *Service) RemoveProjectInformation(uid, pid string) (*entity.Project, error) {

	if !service.conn.SearchLink(uid, pid) {
		return nil, errors.New("unauthorized user, no such link exists")
	}
	if err := service.conn.UnLinkProject(uid, pid); err != nil {
		return nil, err
	}

	project := service.SearchProjectByID(pid)

	for _, value := range service.conn.GetAttachedFiles(pid) {
		if err := service.conn.RemoveFile(value); err != nil {
			return nil, err
		}
	}
	if err := service.conn.RemoveAttachedFiles(pid); err != nil {
		return nil, err
	}

	// -------------------------------------------//
	if err := service.conn.RemoveProject(pid); err != nil {
		return nil, err
	}

	return project, nil
}

// FindProject is a method that is used for searching a projects using search key and other filters.
func (service *Service) FindProject(searchKey string, searchBy string, filterType int64, filterValue1 float64, filterValue2 float64, pageNumber int64) []*entity.Project {

	projects := service.conn.SearchProject(searchKey, searchBy, filterType, filterValue1, filterValue2, pageNumber)
	return projects

}

// GetSentProjects is a method that returns all the projects linked to a user.
func (service *Service) GetSentProjects(uid string) []*entity.Project {

	listOfPids := service.conn.GetLinkedProjects(uid)
	var listOfSentProjects []*entity.Project
	for _, pid := range listOfPids {
		sentProject := service.SearchProjectByID(pid)
		listOfSentProjects = append(listOfSentProjects, sentProject)
	}

	return listOfSentProjects

}

// Validateproject is a method that is used for valdiating project information.
func (service *Service) validateProject(project *entity.Project) entity.ErrorBag {

	var errMap entity.ErrorBag = make(map[string]string)
	if len(project.Title) < 3 {
		errMap["title"] = "Title length too short!"
	}

	if len(project.Description) < 3 {
		errMap["description"] = "Description length too short!"
	}

	if len(project.Details) < 3 {
		errMap["detail"] = "Details length too short!"
	}

	if project.WorkType > 3 {
		errMap["worktype"] = "Unknown worktype!"
	}

	if flag := service.conn.SearchMember("categories", project.Category); !flag {
		errMap["category"] = "Unknown category!"
	}

	if flag := service.conn.SearchMember("subcategories", project.Subcategory); !flag {
		errMap["subcategory"] = "Unknown subcategory!"
	}
	return errMap

}

// CountMember is a method that returns the total number of a member in a table.
func (service *Service) CountMember(tableName string) int {

	members := service.conn.CountMember(tableName)
	return members
}

// AttachFiles is a method that attach a files on certain a given project.
func (service *Service) AttachFiles(pid string, fileName string) error {

	err := service.conn.AttachFiles(pid, fileName)
	return err
}

// SearchMember is a method that checks whether a certain entity is found in a give table.
func (service *Service) SearchMember(tableName string, columnValue string) bool {

	value := service.conn.SearchMember(tableName, columnValue)
	return value
}

// RemoveAttachedFile is a method that removes an attached file from a project.
func (service *Service) RemoveAttachedFile(pid string, fileName string) error {

	err := service.conn.RemoveAttachedFile(pid, fileName)
	return err
}

// RemoveFile is a method that a file from the assets folder.
func (service *Service) RemoveFile(filename string) error {

	err := service.conn.RemoveFile(filename)
	return err
}

// SearchLink is a method that search a connection between a project and a user.
func (service *Service) SearchLink(uid, pid string) bool {

	if !service.conn.SearchLink(uid, pid) {
		return false
	}
	return true
}

// GetCategories return all the categories that are avaliable for project entity.
func (service *Service) GetCategories() []string {
	return service.conn.GetCategories()
}

// GetSubCategories return all the sub categories that are avaliable for project entity.
func (service *Service) GetSubCategories() []string {
	return service.conn.GetSubCategories()
}

// GetSubCategoriesOf return all the sub categories related to a certain category.
func (service *Service) GetSubCategoriesOf(category string) []string {
	return service.conn.GetSubCategoriesOf(category)
}
