package project

import (
	"errors"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
	"github.com/Team-x-AAIT/Freelancing-website/user"
)

// Service is a struct that defines the Project Service type.
type Service struct {
	conn IRepository
}

// IService is an interface that specifies what a Project type can do.
type IService interface {
	PostProject(*entities.Project, string) (string, error)
	SearchProjectByID(string) *entities.Project
	ViewProject(pid string) *user.ProjectUserContainer
	FindProject(string, string, int64, float64, float64, int64) []*entities.Project
	UpdateProject(*entities.Project) (string, error)
	RemoveProjectInformation(string, string) error
	Apply(string, string, string) error
	GetProjectApplicantsID(string) []*ApplicationBag
	GetAppliedFor(string) []*ApplicationBag
	HireApplicant(string, string) error
	RemoveApplicant(string, string) error
	GetSentProjects(uid string) []*entities.Project
	// CheckApplicationStatus(string, string) int64
}

// NewService is a function that returns a new project service type.
func NewService(connection IRepository) IService {
	return &Service{conn: connection}
}

// PostProject is a method that validates and adds a project to the system.
func (service *Service) PostProject(project *entities.Project, uid string) (string, error) {

	if err := service.validateProject(project); err != nil {
		return "", err
	}

	pid, err := service.conn.AddProject(project)
	if err != nil {
		return "", err
	}

	err = service.conn.LinkProject(uid, pid)
	if err != nil {
		return "", err
	}

	return pid, nil
}

// SearchProjectByID is a method that returns projects using the project id.
func (service *Service) SearchProjectByID(pid string) *entities.Project {

	project := service.conn.GetProject(pid)
	attachedFiles := service.conn.GetAttachedFiles(pid)
	project.AttachedFiles = attachedFiles
	return project

}

// ViewProject is a method that returns a project with its owner information.
func (service *Service) ViewProject(pid string) *user.ProjectUserContainer {

	layoutUS := "January 2, 2006"
	projectW := service.SearchProjectByID(pid)
	owner := user.UService.SearchUser(user.URepositoryDB.GetOwner(projectW.ID))

	projectUserContainer := new(user.ProjectUserContainer)
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
func (service *Service) UpdateProject(project *entities.Project) (string, error) {

	if err := service.validateProject(project); err != nil {
		return "", err
	}

	pid, err := service.conn.UpdateProject(project)
	if err != nil {
		return "", err
	}
	return pid, nil
}

// RemoveProjectInformation is a method that is used for removing project and its dependencies.
func (service *Service) RemoveProjectInformation(uid, pid string) error {

	if !service.conn.SearchLink(uid, pid) {
		return errors.New("unauthorized user, no such link exists")
	}
	if err := service.conn.UnLinkProject(uid, pid); err != nil {
		return err
	}
	for _, value := range service.conn.GetAttachedFiles(pid) {
		if err := service.conn.RemoveFile(value); err != nil {
			return err
		}
	}
	if err := service.conn.RemoveAttachedFiles(pid); err != nil {
		return err
	}
	if err := service.conn.RemoveApplicationInfo("", pid); err != nil {
		return err
	}
	if err := service.conn.RemoveProject(pid); err != nil {
		return err
	}

	return nil
}

// FindProject is a method that is used for searching a projects using search key and other filters.
func (service *Service) FindProject(searchKey string, searchBy string, filterType int64, filterValue1 float64, filterValue2 float64, pageNumber int64) []*entities.Project {

	projects := service.conn.SearchProject(searchKey, searchBy, filterType, filterValue1, filterValue2, pageNumber)
	return projects

}

// GetSentProjects is a method that returns all the projects linked to a user.
func (service *Service) GetSentProjects(uid string) []*entities.Project {

	listOfPids := service.conn.GetLinkedProjects(uid)
	var listOfSentProjects []*entities.Project
	for _, pid := range listOfPids {
		sentProject := service.SearchProjectByID(pid)
		listOfSentProjects = append(listOfSentProjects, sentProject)
	}

	return listOfSentProjects

}

// Apply is a method that enables project application process.
func (service *Service) Apply(pid, applicantUID, proposal string) error {

	if service.conn.SearchLink(applicantUID, pid) {
		return errors.New("Can't apply for own project")
	}

	if application := service.conn.GetApplication(pid, applicantUID); application.PID != "" {
		return errors.New("can't apply more than once")
	}

	err := service.conn.AddApplication(pid, applicantUID, proposal)
	if err != nil {
		return err
	}

	err = service.conn.AddApplicationToHistory(pid, applicantUID, proposal)
	if err != nil {
		return err
	}
	return nil
}

// GetAppliedFor is a method that returns all the application the user made based on a certain category.
func (service *Service) GetAppliedFor(uid string) []*ApplicationBag {

	listOfApplications := service.conn.GetUserApplicationHistory(uid)
	for _, value := range listOfApplications {
		project := service.SearchProjectByID(value.PID)
		value.Project = project
		value.Status = service.CheckApplicationStatus(value.PID, value.ApplicantID)
	}

	return listOfApplications

}

// HireApplicant is a method that enables project application acceptance.
func (service *Service) HireApplicant(pid string, applicantUID string) error {

	application := service.conn.GetApplication(pid, applicantUID)
	if application.PID == "" {
		return errors.New("application not found")
	}

	err := service.conn.HireApplicant(pid, applicantUID)
	if err != nil {
		return err
	}

	err = service.conn.UpdateApplicationHistoryTable(pid, applicantUID, application.Proposal, true, false)
	if err != nil {
		return err
	}

	return nil
}

// RemoveApplicant is a method that removes or detaches an applicant from project.
func (service *Service) RemoveApplicant(applicantUID, pid string) error {

	err := service.conn.RemoveApplicationInfo(applicantUID, pid)
	if err != nil {
		return err
	}
	return nil
}

// CheckApplicationStatus is a method that checks users application status and returns its state.
func (service *Service) CheckApplicationStatus(pid, applicantUID string) int64 {

	project := service.SearchProjectByID(pid)
	application := service.conn.GetApplication(pid, applicantUID)
	applicationComparable := service.conn.GetApplicationFromHistory(pid, applicantUID)

	// Means the Project is removed.
	if project.ID == "" {
		// project removed
		return 5
	}
	if application.PID == "" {

		if applicationComparable.PID != "" && applicationComparable.Hired {
			// Fired
			return 4
		}
		if applicationComparable.PID != "" && !applicationComparable.Hired {
			//Applicant removed from project(rejected)
			return 3
		}

	} else {
		if applicationComparable.PID != "" && applicationComparable.Hired {
			// Hired
			return 2
		}
		if applicationComparable.PID != "" && !applicationComparable.Hired {
			// Pending request
			return 1
		}

	}
	return 0
}

// GetProjectApplicantsID is a method that returns a list of users ID and there proposal that applied for a certain project.
func (service *Service) GetProjectApplicantsID(pid string) []*ApplicationBag {

	listOfApplicantsID := service.conn.GetApplicants(pid)
	return listOfApplicantsID
}

// Validateproject is a method that is used for valdiating project information.
func (service *Service) validateProject(project *entities.Project) error {

	switch {
	case len(project.Title) < 3:
		return errors.New("title length too short")
	case len(project.Description) < 3:
		return errors.New("description length too short")
	case len(project.Details) < 3:
		return errors.New("details length too short")
	case project.WorkType > 3:
		return errors.New("unknown worktype")
	}

	if flag := service.conn.SearchMember("categories", project.Category); !flag {
		return errors.New("unknown category")
	}

	if flag := service.conn.SearchMember("subcategories", project.Subcategory); !flag {
		return errors.New("unknown subcategory")
	}
	return nil

}
