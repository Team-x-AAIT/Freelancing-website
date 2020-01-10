package project

import (
	"net/http"
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
