package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Team-x-AAIT/Freelancing-website/stringTools"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// Repository is a struct that define the Repository type.
type Repository struct {
	connection *sql.DB
}

// NewProjectRepository is a function that return new Repository type.
func NewProjectRepository(conn *sql.DB) *Repository {
	return &Repository{connection: conn}
}

// GetProject is a method that searches and deliver the needed project using the project id.
func (psql *Repository) GetProject(pid string) *entity.Project {

	var project entity.Project
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

// AddProject is a method that adds a project to the provided database.
func (psql *Repository) AddProject(project *entity.Project) (string, error) {

	totalNumOfProjects := psql.CountMember("projects")
	totalNumOfProjects++
	project.ID = fmt.Sprintf("PID%d"+stringTools.RandomStringGN(5), totalNumOfProjects)

	for psql.SearchMember("projects", project.ID) {
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

// UpdateProject is a method that updates a project using the project id.
func (psql *Repository) UpdateProject(project *entity.Project) (string, error) {

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

// RemoveProject is a method that remove a project from the database using the project id.
func (psql *Repository) RemoveProject(pid string) error {

	stmt, _ := psql.connection.Prepare(`DELETE FROM projects WHERE id=?`)
	_, err := stmt.Exec(pid)

	if err != nil {
		return err
	}
	return nil
}

// SearchProject is a method that is used for searching a projects from the database using a search-key.
func (psql *Repository) SearchProject(searchKey string, searchBy string, filterType int64, filterValue1 float64, filterValue2 float64, pageNumber int64) []*entity.Project {

	resultLimit := int64(10)
	resultBeginning := pageNumber*resultLimit - resultLimit
	resultEnding := pageNumber*resultLimit - 1
	ctx := context.Background()
	stmt := ""
	rows, err := psql.connection.QueryContext(ctx, stmt)

	switch {
	case searchBy == "title" && filterType != -1:
		stmt = `SELECT * FROM projects WHERE title=? && closed=false
		&& worktype=? && (? < budget < ?) ORDER BY title ASC LIMIT ?,?`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey, filterType,
			filterValue1, filterValue2, resultBeginning, resultEnding)

	case searchBy == "category" && filterType != -1:
		stmt = `SELECT * FROM projects WHERE category = ? && closed=false
		 && worktype = ? && (? < budget < ?) ORDER BY title ASC LIMIT ?,?`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey, filterType,
			filterValue1, filterValue2, resultBeginning, resultEnding)

	case searchBy == "subcategory" && filterType != -1:
		stmt = `SELECT * FROM projects WHERE subcategory=? && closed=false
		&& worktype=? && (? < budget < ?) ORDER BY title ASC LIMIT ?,?`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey, filterType,
			filterValue1, filterValue2, resultBeginning, resultEnding)

	case filterType != -1:
		stmt = `SELECT * FROM projects WHERE (title=? OR category=? OR subcategory=?)
		&& closed=false && worktype=? && (? < budget < ?) ORDER BY title ASC LIMIT ?,?`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey, searchKey,
			searchKey, filterType,
			filterValue1, filterValue2, resultBeginning, resultEnding)

	case searchBy == "title":
		stmt = `SELECT * FROM projects WHERE title = ? && closed=false`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey)

	case searchBy == "category":
		stmt = `SELECT * FROM projects WHERE category = ? && closed=false`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey)
		fmt.Printf("Inside 6")

	case searchBy == "subcategory":
		stmt = `SELECT * FROM projects WHERE subcategory=? && closed=false`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey)

	default:
		stmt = `SELECT * FROM projects WHERE (title=? OR category=? OR subcategory=?)
		&& closed=false ORDER BY title ASC LIMIT ?,?`
		rows, err = psql.connection.QueryContext(ctx, stmt, searchKey, searchKey,
			searchKey, resultBeginning, resultEnding)
	}

	if err != nil {
		return nil
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

		// for postgress is different please please change it don't forget.
		projects = append(projects, project)
	}

	return projects
}

// MarkAsClosed is a method that mark a project as closed by updating the closed filed of a project in the database.
func (psql *Repository) MarkAsClosed(pid string) error {
	stmt, _ := psql.connection.Prepare(`UPDATE projects SET closed=TRUE WHERE pid=?`)
	_, err := stmt.Exec(pid)

	if err != nil {
		return err
	}
	return nil
}

// AttachFiles is a method that adds a filename and its owner that is the project id into attached_files table.
func (psql *Repository) AttachFiles(pid string, fileName string) error {

	stmt, _ := psql.connection.Prepare(`INSERT INTO attached_files (pid, name) VALUES (?,?)`)
	_, err := stmt.Exec(pid, fileName)

	if err != nil {

		panic(err)
		// return err
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

// RemoveAttachedFiles is a method that is used for removing all the files attached to a project from the database.
func (psql *Repository) RemoveAttachedFiles(pid string) error {

	stmt, _ := psql.connection.Prepare(`DELETE FROM attached_files WHERE pid = ?`)
	_, err := stmt.Exec(pid)

	if err != nil {
		return err
	}

	return err

}

// RemoveAttachedFile is a method that is used for removing a file attached to a project from the database.
func (psql *Repository) RemoveAttachedFile(pid string, fileName string) error {

	stmt, _ := psql.connection.Prepare(`DELETE FROM attached_files WHERE pid = ? AND name=?`)
	_, err := stmt.Exec(pid, fileName)

	if err != nil {
		return err
	}

	return err

}

// RemoveFile is a method that removes a given file path from the assets folder.
func (psql *Repository) RemoveFile(filename string) error {

	if err := os.Remove("./assets/" + filename); err != nil {
		return err
	}

	return nil

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

// SearchMember is a method that is used for searching the member of a table where our table name is provided as an argument.
func (psql *Repository) SearchMember(tableName string, columnValue string) bool {

	stmt, _ := psql.connection.Prepare("")
	if tableName == "attached_files" || tableName == "categories" || tableName == "subcategories" {
		stmt, _ = psql.connection.Prepare("SELECT COUNT(*) FROM " + tableName + " WHERE name=?")
	}

	if tableName == "projects" {
		stmt, _ = psql.connection.Prepare("SELECT COUNT(*) FROM " + tableName + " WHERE id=?")
	}

	var totalNumOfMembers int
	row := stmt.QueryRow(columnValue)
	row.Scan(&totalNumOfMembers)

	if totalNumOfMembers > 0 {
		return true
	}
	return false

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

// LinkProject is a method that links a project to user.
func (psql *Repository) LinkProject(uid, pid string) error {

	stmt, _ := psql.connection.Prepare(`INSERT INTO user_project_table (uid, pid) VALUES (?,?)`)
	_, err := stmt.Exec(uid, pid)

	if err != nil {
		return err
	}
	return nil

}

// UnLinkProject is a method that a unlinks a project from a user.
func (psql *Repository) UnLinkProject(uid, pid string) error {

	stmt, _ := psql.connection.Prepare(`DELETE FROM user_project_table WHERE uid=? && pid=?`)
	_, err := stmt.Exec(uid, pid)

	if err != nil {
		return err
	}
	return nil

}

// SearchLink is a method that is used for identifying a certain link exists between a project and a user.
func (psql *Repository) SearchLink(uid, pid string) bool {

	stmt, _ := psql.connection.Prepare("SELECT COUNT(*) FROM user_project_table WHERE uid=? && pid=?")

	var totalNumOfMembers int
	row := stmt.QueryRow(uid, pid)
	row.Scan(&totalNumOfMembers)

	if totalNumOfMembers > 0 {
		return true
	}
	return false
}

// GetCategories retrives all the categories stored in the database.
func (psql *Repository) GetCategories() []string {

	var category string
	var listOfCategories []string
	rows, _ := psql.connection.Query("SELECT name FROM categories")
	if rows.Next() {
		rows.Scan(&category)
		listOfCategories = append(listOfCategories, category)
	}
	return listOfCategories
}

// GetSubCategories retrives all the sub categories stored in the database.
func (psql *Repository) GetSubCategories() []string {

	var subcCategory string
	var listOfSubcCategory []string
	stmt, err := psql.connection.Prepare("SELECT name FROM subcategories")
	if err != nil {
		return nil
	}
	rows, _ := stmt.Query()
	if rows.Next() {
		rows.Scan(&subcCategory)
		listOfSubcCategory = append(listOfSubcCategory, subcCategory)

	}
	return listOfSubcCategory
}

// GetCategoryID retrives the category id of a certain category.
func (psql *Repository) GetCategoryID(category string) int {

	var categoryID int
	stmt, err := psql.connection.Prepare("SELECT id FROM categories WHERE name=?")
	if err != nil {
		return -1
	}
	row := stmt.QueryRow(category)
	row.Scan(&categoryID)
	return categoryID
}

// GetSubCategoriesOf retrives the sub categories of a given category.
func (psql *Repository) GetSubCategoriesOf(category string) []string {

	categoryID := psql.GetCategoryID(category)

	var subcCategory string
	var listOfSubcCategory []string
	stmt, err := psql.connection.Prepare("SELECT name FROM subcategories WHERE cid=?")
	if err != nil {
		return nil
	}
	rows, _ := stmt.Query(categoryID)
	if rows.Next() {
		rows.Scan(&subcCategory)
		listOfSubcCategory = append(listOfSubcCategory, subcCategory)

	}
	return listOfSubcCategory
}
