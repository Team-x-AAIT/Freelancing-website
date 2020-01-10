package project

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
	"github.com/Team-x-AAIT/Freelancing-website/user"
)

var (
	// PRepositoryDB is a pointer to the user PsqlProjectRepository type.
	PRepositoryDB IRepository
	// PService is a pointer to the project Service type.
	PService IService
)

// Bag is used to store both the current user and project information of a selected project.
type Bag struct {
	CurrentUser          *entities.User
	ProjectUserContainer *user.ProjectUserContainer
}

// Applicants is a struct that holds the use information and his/her proposal statment.
type Applicants struct {
	ApplicantUID string
	Firstname    string
	Lastname     string
	Email        string
	Gender       string
	JobTitle     string
	PhoneNumber  string
	ProfilePic   string
	Rating       float32
	Proposal     string
	Hired        bool
}

// PostProject is a Handler func that initaite project posting process.
func PostProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	uid := loggedInUser.UID

	title := r.FormValue("title")
	description := r.FormValue("description")
	details := r.FormValue("details")
	category := r.FormValue("category")
	subcategory := r.FormValue("subcategory")
	budgetS := r.FormValue("budget")
	worktypeS := r.FormValue("worktype")

	budget, err := strconv.ParseFloat(budgetS, 0)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	worktype, err := strconv.ParseInt(worktypeS, 0, 0)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	projectHolder := entities.NewProject(title, description, details, category, subcategory, budget, worktype)
	pid, err := PService.PostProject(projectHolder, uid)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	ResourceExtractor(pid, "files", "project", r)
	user.Temp.ExecuteTemplate(w, "Dashboard.html", loggedInUser)
	//Don't forget to change it to ajax request.
	// w.Write([]byte("okay"))

}

// OnCreateProject is a Handler func that display a new page for creating a new project.
func OnCreateProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}
	user.Temp.ExecuteTemplate(w, "PostProject.html", loggedInUser)
}

// ReviewProject is a Handler func that display the created project for editing.
func ReviewProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	pid := r.FormValue("pid")
	project := PService.SearchProjectByID(pid)
	user.Temp.ExecuteTemplate(w, "ReviewProject.html", project)

}

// UpdateProject is a Handler func that initaite project updating process.
func UpdateProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	pid := r.FormValue("pid")
	title := r.FormValue("title")
	description := r.FormValue("description")
	details := r.FormValue("details")
	category := r.FormValue("category")
	subcategory := r.FormValue("subcategory")
	budgetS := r.FormValue("budget")
	worktypeS := r.FormValue("worktype")
	// Is a string that contain the names of file that are needed to removed.
	attachedFilesR := r.FormValue("attachedFilesR")

	budget, err := strconv.ParseFloat(budgetS, 0)
	if err != nil {
		panic(err)
	}
	worktype, err := strconv.ParseInt(worktypeS, 0, 0)
	if err != nil {
		panic(err)
	}

	projectHolder := entities.NewProject(title, description, details, category, subcategory, budget, worktype)

	projectHolder.ID = pid
	pid, err = PService.UpdateProject(projectHolder)
	if err != nil {
		panic(err)
	}

	ResourceExtractor(pid, "files", "project", r)

	var beginIndex int
	for index, char1 := range attachedFilesR {
		if char1 == ',' {
			filename := attachedFilesR[beginIndex:index]
			if err = RemoveAsset(pid, filename); err != nil {
				panic(err)
			}
			beginIndex = index + 1
		}
	}

}

// RemoveProject is a Handler func that initaite project Removing process.
func RemoveProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	uid := loggedInUser.UID
	pid := r.FormValue("pid")

	if err := PService.RemoveProjectInformation(uid, pid); err != nil {
		panic(err)
	}

}
