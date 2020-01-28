package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// Repository is a struct that define the Repository type.
type Repository struct {
	connection *sql.DB
}

// NewUserRepository is a function that return new IRepository type.
func NewUserRepository(conn *sql.DB) *Repository {
	return &Repository{connection: conn}
}

// AddUser is a method that adds a user to the provided database.
func (psql *Repository) AddUser(user *entity.User) error {

	totalNumOfUsers := psql.CountMember("users")

	user.UID = fmt.Sprintf("UID%d", totalNumOfUsers+1)

	for psql.SearchMember("users", user.UID) {
		totalNumOfUsers++
		user.UID = fmt.Sprintf("UID%d", totalNumOfUsers)
	}

	stmt, _ := psql.connection.Prepare(`INSERT INTO Users (uid,first_name, last_name, password, phonenumber, email, 
		job_title, country, city, gender,cv, profile_pic, bio, prefe, rating) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
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
		user.Prefe,
		user.Rating)

	if err != nil {
		return err
	}
	return nil
}

// SearchUser is a method that search for a User type from the 'users' table in the database using the provided Email address.
func (psql *Repository) SearchUser(identifier string) *entity.User {

	stmt, err := psql.connection.Prepare("SELECT * FROM Users WHERE email=? || uid=?")
	row := stmt.QueryRow(identifier, identifier)
	if err != nil {
		panic(err)
	}
	var user entity.User
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
		&user.Prefe,
		&user.Rating)

	return &user
}

// UpdateUser is a method that updates the user data in the database using the provided parameters.
func (psql *Repository) UpdateUser(user *entity.User) error {
	stmt, _ := psql.connection.Prepare(`UPDATE Users SET first_name = ?, last_name = ?, phonenumber = ?, email = ?, job_title = ?,
		country = ?, city = ?, gender = ?, cv = ?, profile_pic = ?, bio = ?, prefe = ? WHERE uid = ?`)
	_, err := stmt.Exec(
		user.Firstname,
		user.Lastname,
		user.Phonenumber,
		user.Email,
		user.JobTitle,
		user.Country,
		user.City,
		user.Gender,
		user.CV,
		user.ProfilePic,
		user.Bio,
		user.Prefe,
		user.UID)

	if err != nil {
		return err
	}
	return nil

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

// UpdateTPUsers is a method that updates the third party user data in the database using the provided parameters.
func (psql *Repository) UpdateTPUsers(uid, email string) error {
	stmt, _ := psql.connection.Prepare(`UPDATE tp_users SET email=? WHERE uid=?`)
	_, err := stmt.Exec(email, uid)

	if err != nil {
		return err
	}
	return nil
}

// SearchTPUser is a method that search for a User type from the 'tp_users' table in database using the provided Email address.
func (psql *Repository) SearchTPUser(identifier string) (*entity.User, string) {

	stmt, err := psql.connection.Prepare("SELECT * FROM tp_users WHERE email=? || uid=?")
	row := stmt.QueryRow(identifier, identifier)
	if err != nil {
		panic(err)
	}
	var from string
	var user entity.User
	row.Scan(
		&user.UID,
		&user.Email,
		&user.Password,
		&from)
	return &user, from
}

// RemoveUser remove a user from user table.
func (psql *Repository) RemoveUser(uid string) (*entity.User, error) {

	userHolder := psql.SearchUser(uid)

	if userHolder.UID == "" {
		return nil, errors.New("user not found")
	}

	stmt, _ := psql.connection.Prepare(`DELETE FROM users WHERE uid=?`)
	_, err := stmt.Exec(uid)

	if err != nil {
		return nil, err
	}
	return userHolder, nil

}

// RemoveTPUser remove a user from user_tp table.
func (psql *Repository) RemoveTPUser(uid string) (*entity.User, error) {

	userHolder, _ := psql.SearchTPUser(uid)

	if userHolder.UID == "" {
		return nil, errors.New("user not found")
	}

	stmt, _ := psql.connection.Prepare(`DELETE FROM tp_users WHERE uid=?`)
	_, err := stmt.Exec(uid)

	if err != nil {
		return nil, err
	}
	return userHolder, nil
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

// AddMatchTag is a method that adds a matchTag to the match tag table.
func (psql *Repository) AddMatchTag(uid string, category string, subcategory string, worktype int) error {

	stmt, _ := psql.connection.Prepare(`INSERT INTO match_tag_table (uid,category,subcategory,worktype
		) VALUES (?,?,?,?)`)
	_, err := stmt.Exec(
		uid,
		category,
		subcategory,
		worktype,
	)

	if err != nil {
		panic(err)
	}
	return nil
}

// RemoveMatchTag is a method that removes a matchTag from the match tag table.
func (psql *Repository) RemoveMatchTag(uid string, category string, subcategory string, worktype int) error {

	stmt, _ := psql.connection.Prepare(`DELETE FROM match_tag_table WHERE uid=? && category=? && 
	subcategory=? && worktype=?`)
	_, err := stmt.Exec(
		uid,
		category,
		subcategory,
		worktype,
	)

	if err != nil {
		panic(err)
	}
	return nil
}

// GetUserMatchTags is a method that returns all the match tag the user defined.
func (psql *Repository) GetUserMatchTags(uid string) []*entity.MatchTag {

	stmt, err := psql.connection.Prepare("SELECT * FROM match_tag_table WHERE uid=?")
	rows, err := stmt.Query(uid)
	if err != nil {
		panic(err)
	}

	var matchTagStrore []*entity.MatchTag
	for rows.Next() {

		matchTag := new(entity.MatchTag)

		rows.Scan(
			&matchTag.UID,
			&matchTag.Category,
			&matchTag.Subcategory,
			&matchTag.WorkType)
		matchTagStrore = append(matchTagStrore, matchTag)
	}

	return matchTagStrore
}

// SearchProjectWMatchTag is a method that returns all the projects that match the provided tag.
func (psql *Repository) SearchProjectWMatchTag(matchTag *entity.MatchTag) []*entity.Project {

	query := `SELECT * FROM projects WHERE`
	if matchTag.Category != "" {
		query += fmt.Sprintf(` category='%s'`, matchTag.Category)
	}
	if matchTag.Subcategory != "" {
		query += fmt.Sprintf(` && subcategory='%s'`, matchTag.Subcategory)
	}
	if matchTag.WorkType != 4 {
		query += fmt.Sprintf(` && worktype=%d`, matchTag.WorkType)
	}

	stmt, err := psql.connection.Prepare(query)
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	var projects []*entity.Project

	for rows.Next() {
		project := new(entity.Project)
		err = rows.Scan(&project.ID, &project.Title,
			&project.Description, &project.Details,
			&project.Category, &project.Subcategory,
			&project.Budget, &project.WorkType,
			&project.Closed, &project.CreatedAt)

		if err != nil {
			panic(err)
		}

		projects = append(projects, project)
	}

	return projects

}

// RemoveFileDB is a method that is used for clearing a certain column.
func (psql *Repository) RemoveFileDB(uid string, columnName string) string {

	ctx := context.Background()
	var value string

	stmt := `SELECT ` + columnName + ` FROM users WHERE uid = ?`
	row := psql.connection.QueryRowContext(ctx, stmt, uid)

	row.Scan(&value)

	stmt = `UPDATE users SET ` + columnName + `="" WHERE uid = ?`
	_, err := psql.connection.QueryContext(ctx, stmt, uid)

	if err != nil {
		return ""
	}

	return value

}

// RemoveFile is a method that removes a given file path from the assets folder.
func (psql *Repository) RemoveFile(filename string, folder string) error {

	if err := os.Remove("../../ui/assets/" + folder + "/" + filename); err != nil {
		return err
	}
	return nil
}

// SearchMember is a method that is used for searching the member of a table where our table name is provided as an argument.
func (psql *Repository) SearchMember(tableName string, columnValue string) bool {

	stmt, _ := psql.connection.Prepare("")
	if tableName == "attached_files" || tableName == "categories" || tableName == "subcategories" {
		stmt, _ = psql.connection.Prepare("SELECT COUNT(*) FROM " + tableName + " WHERE name=?")
	}

	if tableName == "projects" {
		stmt, _ = psql.connection.Prepare("SELECT COUNT(*) FROM " + tableName + " WHERE id=?")
	}

	if tableName == "users" {
		stmt, _ = psql.connection.Prepare("SELECT COUNT(*) FROM " + tableName + " WHERE uid=?")
	}

	var totalNumOfMembers int
	row := stmt.QueryRow(columnValue)
	row.Scan(&totalNumOfMembers)

	if totalNumOfMembers > 0 {
		return true
	}
	return false

}

// GetOwner is a method that returns the owner of a project by searching through the project user table.
func (psql *Repository) GetOwner(pid string) string {
	stmt, _ := psql.connection.Prepare("SELECT uid FROM user_project_table WHERE pid=?")

	var uid string
	row := stmt.QueryRow(pid)
	row.Scan(&uid)

	return uid
}
