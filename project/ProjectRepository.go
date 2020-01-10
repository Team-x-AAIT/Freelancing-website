package project

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
	AttachFiles(string, string) error
	GetProject(string) *entities.Project
	GetLinkedProjects(string) []string
	UpdateProject(*entities.Project) (string, error)
	RemoveProject(string) error
	LinkProject(string, string) error
	UnLinkProject(string, string) error
	SearchLink(string, string) bool
	SearchProject(string, string, int64, float64, float64, int64) []*entities.Project
	AddApplication(string, string, string) error
	AddApplicationToHistory(string, string, string) error
	GetApplication(string, string) *ApplicationBag
	GetApplicationFromHistory(string, string) *ApplicationBag
	GetUserApplicationHistory(string) []*ApplicationBag
	GetAttachedFiles(string) []string
	RemoveAttachedFiles(string) error
	RemoveAttachedFile(string, string) error
	RemoveFile(string) error
	RemoveApplicationInfo(string, string) error
	GetApplicants(string) []*ApplicationBag
	HireApplicant(string, string) error
	UpdateApplicationTable(string, string, string, bool) error
	UpdateApplicationHistoryTable(string, string, string, bool, bool) error
	RemoveUnHiredApplicants(string) error
	MarkAsClosed(string) error
}

// NewRepository is a function that return new Repository type.
func NewRepository(conn *sql.DB) IRepository {
	return &Repository{connection: conn}
}

// AddProject is a method that adds a project to the provided database.
func (psql *Repository) AddProject(project *entities.Project) (string, error) {

	totalNumOfProjects := psql.CountMember("projects")
	totalNumOfProjects++
	project.ID = fmt.Sprintf("PID%d", totalNumOfProjects)

	for PRepositoryDB.SearchMember("projects", project.ID) {
		totalNumOfProjects++
		project.ID = fmt.Sprintf("PID%d", totalNumOfProjects)
	}

	stmt, _ := psql.connection.Prepare(`INSERT INTO projects (id, title, description, details, category, subcategory, budget, worktype, created_at)
	 VALUES (?,?,?,?,?,?,?,?,?)`)
	_, err := stmt.Exec(
		project.ID,
		project.Title,
		project.Description,
		project.Details,
		project.Category,
		project.Subcategory,
		project.Budget,
		project.WorkType,
		time.Now())

	if err != nil {
		return "", err
	}
	return project.ID, nil
}

// AttachFiles is a method that adds a filename and its owner that is the project id into attached_files table.
func (psql *Repository) AttachFiles(pid string, fileName string) error {

	stmt, _ := psql.connection.Prepare(`INSERT INTO attached_files (pid, name) VALUES (?,?)`)
	_, err := stmt.Exec(pid, fileName)

	if err != nil {
		return err
	}
	return nil
}

// GetAttachedFiles is a method that returns names of the files that are attached to the project.
func (psql *Repository) GetAttachedFiles(pid string) []string {

	var name string
	var attachedFiles []string
	stmt, _ := psql.connection.Prepare(`SELECT name FROM attached_files WHERE pid=?`)
	rows, err := stmt.Query(pid)

	if err != nil {
		return attachedFiles
	}

	for rows.Next() {
		rows.Scan(&name)
		attachedFiles = append(attachedFiles, name)
	}

	return attachedFiles
}

// GetProject is a method that searches and deliver the needed project using the project id.
func (psql *Repository) GetProject(pid string) *entities.Project {

	var project entities.Project
	stmt, err := psql.connection.Prepare("SELECT * FROM projects WHERE id=?")
	if err != nil {
		return nil
	}
	row := stmt.QueryRow(pid)
	row.Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		&project.Details,
		&project.Category,
		&project.Subcategory,
		&project.Budget,
		&project.WorkType,
		&project.Closed,
		&project.CreatedAt,
	)

	return &project

}

// UpdateProject is a method that updates a project using the project id.
func (psql *Repository) UpdateProject(project *entities.Project) (string, error) {

	stmt, _ := psql.connection.Prepare(`UPDATE projects SET title=?, description=?, 
	details=?, category=?, 
	subcategory=?, budget=?, 
	worktype=? WHERE id = ?`)
	_, err := stmt.Exec(
		project.Title,
		project.Description,
		project.Details,
		project.Category,
		project.Subcategory,
		project.Budget,
		project.WorkType,
		project.ID)

	if err != nil {
		return "", err
	}
	return project.ID, nil
}

// GetLinkedProjects is a method that returns all the linked projects id to a certain user.
func (psql *Repository) GetLinkedProjects(uid string) []string {

	var listOfPid []string
	stmt, err := psql.connection.Prepare("SELECT pid FROM user_project_table WHERE uid=?")
	if err != nil {
		return nil
	}
	rows, _ := stmt.Query(uid)
	for rows.Next() {
		var pid string
		rows.Scan(&pid)
		listOfPid = append(listOfPid, pid)
	}

	return listOfPid
}

// AddApplication is a method that adds a user application for a project to a database.
func (psql *Repository) AddApplication(pid, applicantUID, proposal string) error {

	stmt, _ := psql.connection.Prepare(`INSERT INTO application_table (pid, applicant_uid, proposal) 
	VALUES (?,?,?)`)
	_, err := stmt.Exec(pid, applicantUID, proposal)

	if err != nil {
		return err
	}
	return nil

}

// AddApplicationToHistory is a method that adds a user application history for a project to a database.
func (psql *Repository) AddApplicationToHistory(pid, applicantUID, proposal string) error {

	stmt, _ := psql.connection.Prepare(`INSERT INTO application_history_table (pid, applicant_uid, proposal, applied_at) 
	VALUES (?,?,?,?)`)
	_, err := stmt.Exec(pid, applicantUID, proposal, time.Now())

	if err != nil {
		return err
	}
	return nil

}

// GetApplication is a method that checks weather a certain application is
// avaliable in application table and returns the application.
func (psql *Repository) GetApplication(pid, uid string) *ApplicationBag {

	stmt, _ := psql.connection.Prepare("SELECT * FROM application_table WHERE applicant_uid=? && pid=?")

	row := stmt.QueryRow(uid, pid)
	applicationBag := new(ApplicationBag)
	row.Scan(&applicationBag.PID, &applicationBag.ApplicantID, &applicationBag.Proposal, &applicationBag.Hired)

	return applicationBag

}

// GetApplicationFromHistory is a method that checks weather a certain application is
// avaliable in application history table and returns the application.
func (psql *Repository) GetApplicationFromHistory(pid, uid string) *ApplicationBag {

	stmt, _ := psql.connection.Prepare("SELECT * FROM application_history_table WHERE applicant_uid=? && pid=?")

	row := stmt.QueryRow(uid, pid)
	applicationBag := new(ApplicationBag)
	row.Scan(&applicationBag.PID, &applicationBag.ApplicantID,
		&applicationBag.Proposal, &applicationBag.Hired,
		&applicationBag.Seen, &applicationBag.CreatedAt)

	return applicationBag
}

// GetUserApplicationHistory is a method that returns list of application that corresponds to a user id.
func (psql *Repository) GetUserApplicationHistory(uid string) []*ApplicationBag {

	stmt, _ := psql.connection.Prepare("SELECT * FROM application_history_table WHERE applicant_uid=?")

	var listOfApplication []*ApplicationBag
	rows, err := stmt.Query(uid)

	if err != nil {
		return nil
	}

	for rows.Next() {
		applicationBag := new(ApplicationBag)
		rows.Scan(&applicationBag.PID, &applicationBag.ApplicantID,
			&applicationBag.Proposal, &applicationBag.Hired,
			&applicationBag.Seen, &applicationBag.CreatedAt)
		listOfApplication = append(listOfApplication, applicationBag)
	}

	return listOfApplication
}

// HireApplicant is a method that changes the state of an applicant in application_table.
func (psql *Repository) HireApplicant(pid string, applicantUID string) error {

	stmt, _ := psql.connection.Prepare(`UPDATE application_table SET hired=TRUE WHERE pid=? && applicant_uid=?`)
	_, err := stmt.Exec(pid, applicantUID)

	if err != nil {
		return err
	}
	return nil

}
