package service

import (
	"errors"

	"github.com/Team-x-AAIT/Freelancing-website/application"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// Service is a struct that defines the application Service type.
type Service struct {
	conn application.IRepository
}

// NewApplicationService is a function that returns a new application service type.
func NewApplicationService(connection application.IRepository) *Service {
	return &Service{conn: connection}
}

// Apply is a method that enables project application process.
func (service *Service) Apply(pid, applicantUID, proposal string) error {

	if application, err := service.conn.GetApplication(pid, applicantUID); application.PID != "" || err != nil {
		return errors.New("can't apply more than once")
	}

	previouseApplication := service.GetUserApplicationHistory(applicantUID)
	for _, pApplication := range previouseApplication {
		if pApplication.PID == pid {
			return errors.New("can't apply more than once")
		}
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

// GetUserApplicationHistory is a method that returns the application history of a user.
func (service *Service) GetUserApplicationHistory(uid string) []*entity.ApplicationBag {

	listOfApplications := service.conn.GetUserApplicationHistory(uid)
	return listOfApplications
}

// GetApplication is a method that returns an application of a user for a certain project.
func (service *Service) GetApplication(applicantUID, pid string) (*entity.ApplicationBag, error) {

	application, err := service.conn.GetApplication(applicantUID, pid)

	if err != nil {
		return nil, err
	}
	if application.PID == "" {
		return nil, errors.New("application not found")
	}
	return application, err
}

// GetProjectApplicantsID is a method that returns a list of users ID and there proposal that applied for a certain project.
func (service *Service) GetProjectApplicantsID(pid string) []*entity.ApplicationBag {

	listOfApplicantsID := service.conn.GetApplicants(pid)
	return listOfApplicantsID
}

// HireApplicant is a method that enables project application acceptance.
func (service *Service) HireApplicant(pid string, applicantUID string) error {

	application, err := service.conn.GetApplication(pid, applicantUID)
	if application.PID == "" || err != nil {
		return errors.New("application not found")
	}

	err = service.conn.HireApplicant(pid, applicantUID)
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
func (service *Service) RemoveApplicant(applicantUID, pid string) (*entity.ApplicationBag, error) {

	application, err := service.conn.GetApplication(applicantUID, pid)
	if err != nil {
		return nil, err
	}
	err = service.conn.RemoveApplicationInfo(applicantUID, pid)
	if err != nil {
		return nil, err
	}
	return application, nil
}

// RemoveApplicationInfo is a method that removes all the application informations of a particular application.
func (service *Service) RemoveApplicationInfo(applicantUID, pid string) error {

	err := service.conn.RemoveApplicationInfo(applicantUID, pid)
	if err != nil {
		return err
	}
	return nil
}

// CheckApplicationStatus is a method that checks users application status and returns its state.
func (service *Service) CheckApplicationStatus(project *entity.Project, applicantUID string) int64 {

	application, _ := service.conn.GetApplication(applicantUID, project.ID)
	applicationComparable := service.conn.GetApplicationFromHistory(project.ID, applicantUID)

	// Means the Project is removed.
	if project.ID == "" {
		// project removed
		return 4
	}
	if application.PID == "" {

		if applicationComparable.PID != "" && applicationComparable.Hired {
			// Fired
			return 3
		}
		if applicationComparable.PID != "" && !applicationComparable.Hired {
			//Applicant removed from project(rejected)
			return 2
		}

	} else {
		if applicationComparable.PID != "" && applicationComparable.Hired {
			// Hired
			return 1
		}
		if applicationComparable.PID != "" && !applicationComparable.Hired {
			// Pending request
			return 0
		}

	}
	return 0
}

// GetAppliedFor is a method that returns all the application the user made based on a certain category.
// func (service *Service) GetAppliedFor(uid string) []*entity.ApplicationBag {

// 	listOfApplications := service.conn.GetUserApplicationHistory(uid)
// 	for _, value := range listOfApplications {
// 		project := service.SearchProjectByID(value.PID)
// 		value.Project = project
// 		value.Status = service.CheckApplicationStatus(project, value.ApplicantID)
// 	}

// 	return listOfApplications

// }
