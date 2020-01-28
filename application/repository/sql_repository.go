package repository

import (
	"database/sql"
	"time"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// Repository is a struct that define the Repository type.
type Repository struct {
	connection *sql.DB
}

// NewApplicationRepository is a function that return new Repository type.
func NewApplicationRepository(conn *sql.DB) *Repository {
	return &Repository{connection: conn}
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

// GetApplicants is a method that returns all the applicants id's attached to a project.
func (psql *Repository) GetApplicants(pid string) []*entity.ApplicationBag {

	stmt, _ := psql.connection.Prepare("SELECT applicant_uid, proposal, hired FROM application_table WHERE pid=?")

	var listOfApplication []*entity.ApplicationBag
	rows, err := stmt.Query(pid)

	if err != nil {
		return nil
	}

	for rows.Next() {
		applicationBag := new(entity.ApplicationBag)
		rows.Scan(&applicationBag.ApplicantID, &applicationBag.Proposal, &applicationBag.Hired)
		listOfApplication = append(listOfApplication, applicationBag)
	}

	return listOfApplication

}

// GetApplication is a method that checks weather a certain application is
// avaliable in application table and returns the application.
func (psql *Repository) GetApplication(uid, pid string) (*entity.ApplicationBag, error) {

	stmt, err := psql.connection.Prepare("SELECT * FROM application_table WHERE applicant_uid=? && pid=?")

	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(uid, pid)
	applicationBag := new(entity.ApplicationBag)
	row.Scan(&applicationBag.PID, &applicationBag.ApplicantID, &applicationBag.Proposal, &applicationBag.Hired)

	return applicationBag, nil

}

// GetApplicationFromHistory is a method that checks weather a certain application is
// avaliable in application history table and returns the application.
func (psql *Repository) GetApplicationFromHistory(pid, uid string) *entity.ApplicationBag {

	stmt, _ := psql.connection.Prepare("SELECT * FROM application_history_table WHERE applicant_uid=? && pid=?")

	row := stmt.QueryRow(uid, pid)
	applicationBag := new(entity.ApplicationBag)
	row.Scan(&applicationBag.PID, &applicationBag.ApplicantID,
		&applicationBag.Proposal, &applicationBag.Hired,
		&applicationBag.Seen, &applicationBag.CreatedAt)

	return applicationBag
}

// GetUserApplicationHistory is a method that returns list of application that corresponds to a user id.
func (psql *Repository) GetUserApplicationHistory(uid string) []*entity.ApplicationBag {

	stmt, _ := psql.connection.Prepare("SELECT * FROM application_history_table WHERE applicant_uid=?")

	var listOfApplication []*entity.ApplicationBag
	rows, err := stmt.Query(uid)

	if err != nil {
		return nil
	}

	for rows.Next() {
		applicationBag := new(entity.ApplicationBag)
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

// UpdateApplicationTable is a method that is used for updating application table.
func (psql *Repository) UpdateApplicationTable(pid string, applicantUID string, proposal string, status bool) error {

	stmt, _ := psql.connection.Prepare(`UPDATE application_table SET proposal=?, hired=? 
	WHERE applicant_uid=? && pid=?`)
	_, err := stmt.Exec(proposal, status, applicantUID, pid)

	if err != nil {
		return err
	}
	return nil
}

// UpdateApplicationHistoryTable is a method that is used for updating application history table.
func (psql *Repository) UpdateApplicationHistoryTable(pid string, applicantUID string, proposal string, status, seen bool) error {

	stmt, _ := psql.connection.Prepare(`UPDATE application_history_table SET proposal=?, hired=?, seen=?
	WHERE applicant_uid=? && pid=?`)

	_, err := stmt.Exec(proposal, status, seen, applicantUID, pid)

	if err != nil {
		return err
	}
	return nil
}

// RemoveUnHiredApplicants is a method that removes all unhired applicants and unlink them from the project.
func (psql *Repository) RemoveUnHiredApplicants(pid string) error {

	stmt, _ := psql.connection.Prepare(`DELETE FROM application_table WHERE pid=? && hire=FALSE`)
	_, err := stmt.Exec(pid)

	if err != nil {
		return err
	}
	return nil
}

// RemoveApplicationInfo is a method that remove all the applicants attached to a project.
func (psql *Repository) RemoveApplicationInfo(uid, pid string) error {

	var err error
	// Means just removing the application linked with a certain project.
	if uid == "" {
		stmt, _ := psql.connection.Prepare(`DELETE FROM application_table WHERE pid = ?`)
		_, err = stmt.Exec(pid)

		// Removing an application related to a certain applicant id.
	} else {
		stmt, _ := psql.connection.Prepare(`DELETE FROM application_table WHERE applicant_uid=? && pid = ?`)
		_, err = stmt.Exec(uid, pid)
	}

	if err != nil {
		return err
	}
	return err

}
