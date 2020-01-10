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

// ViewProject is a Handler func that enables applicant to view posted project.
func ViewProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	pid := r.FormValue("pid")
	projectUserContainer := PService.ViewProject(pid)
	value := Bag{CurrentUser: loggedInUser, ProjectUserContainer: projectUserContainer}
	user.Temp.ExecuteTemplate(w, "ViewProject.html", value)

}

// GetSentProjects is a Handler func that sends all the projects linked to a user.
func GetSentProjects(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}
	listOfSentProjects := PService.GetSentProjects(loggedInUser.UID)
	json, err := json.Marshal(listOfSentProjects)
	if err != nil {
		panic(err)
	}

	w.Write(json)
}

// SearchProject is a Handler func that initaite projet searching process.
func SearchProject(w http.ResponseWriter, r *http.Request) {
	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	searchKey := r.FormValue("searchKey")
	searchBy := r.FormValue("searchBy")
	filterTypeS := r.FormValue("filterType")
	filterValue1S := r.FormValue("filterValue1")
	filterValue2S := r.FormValue("filterValue2")
	pageNumS := r.FormValue("pageNum")

	filterType, err := strconv.ParseInt(filterTypeS, 0, 0)
	if err != nil {
		filterType = -1
	}
	filterValue1, err := strconv.ParseFloat(filterValue1S, 0)
	if err != nil {
		filterType = 0
	}
	filterValue2, err := strconv.ParseFloat(filterValue2S, 0)
	if err != nil {
		filterType = 0
	}
	pageNum, err := strconv.ParseInt(pageNumS, 0, 0)
	if err != nil {
		filterType = 1
	}

	projects := PService.FindProject(searchKey, searchBy, filterType, filterValue1, filterValue2, pageNum)

	temp := template.Must(template.ParseFiles("./viewProjects.html"))
	temp.Execute(w, projects)

}

// ViewApplicationRequest is a Handler func that enables user to view all the applicants to a certain project.
func ViewApplicationRequest(w http.ResponseWriter, r *http.Request) {

	pid := r.FormValue("pid")
	listOfApplicantsID := PService.GetProjectApplicantsID(pid)
	var listOfProjectApplicants []Applicants

	for _, value := range listOfApplicantsID {
		applicantInfo := user.UService.SearchUser(value.ApplicantID)
		projectApplicant := Applicants{
			ApplicantUID: applicantInfo.UID,
			Firstname:    applicantInfo.Firstname,
			Lastname:     applicantInfo.Lastname,
			Email:        applicantInfo.Email,
			Gender:       applicantInfo.Gender,
			JobTitle:     applicantInfo.JobTitle,
			PhoneNumber:  applicantInfo.Phonenumber,
			ProfilePic:   applicantInfo.ProfilePic,
			Rating:       applicantInfo.Rating,
			Proposal:     value.Proposal,
			Hired:        value.Hired}

		listOfProjectApplicants = append(listOfProjectApplicants, projectApplicant)
	}
	temp := template.Must(template.ParseFiles("./viewApplicants.html"))
	temp.Execute(w, listOfProjectApplicants)

}

// ApplyForProject is a Handler func that initaite user application process.
func ApplyForProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := user.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login_Page", http.StatusSeeOther)
		return
	}

	uid := loggedInUser.UID
	pid := r.FormValue("pid")
	proposal := r.FormValue("proposal")

	if err := PService.Apply(pid, uid, proposal); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("okay"))
	}

}

// HireApplicant is a Handler funct that initaite client hiring processe.
func HireApplicant(w http.ResponseWriter, r *http.Request) {

	pid := r.FormValue("pid")
	applicantUID := r.FormValue("applicant_UID")

	err := PService.HireApplicant(pid, applicantUID)
	if err != nil {
		panic(err)
	}
}

// RemoveApplication is a Handler func that initaite application removal process
func RemoveApplication(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("pid")
	applicantUID := r.FormValue("applicant_UID")

	err := PService.RemoveApplicant(applicantUID, pid)
	if err != nil {
		panic(err)
	}
}

// ViewAsset is used for diplaying assets.
func ViewAsset(w http.ResponseWriter, r *http.Request) {
	asset := r.FormValue("asset")
	folder := r.FormValue("type")
	if folder == "cv" {
		asset = "./assets/cv/" + asset
	}
	if folder == "attached_file" {
		asset = "./assets/attached_files/" + asset
	}

	temp := template.Must(template.ParseFiles("./assetViewer.html"))
	temp.Execute(w, asset)
}

// RemoveAsset is a function that removes a given file path from the system.
func RemoveAsset(pid string, filename string) error {

	if err := PRepositoryDB.RemoveAttachedFile(pid, filename); err != nil {
		return err
	}
	if err := PRepositoryDB.RemoveFile(filename); err != nil {
		return err
	}

	return nil

}

// ResourceExtractor is a function that extract file from a request.
func ResourceExtractor(pid string, name string, fileType string, r *http.Request) {

	err := r.ParseMultipartForm(200000)
	if err != nil {
		panic(err)
	}
	formdata := r.MultipartForm
	files := formdata.File[name]

	for _, file := range files {
		tempFile, err := file.Open()
		defer tempFile.Close()
		if err != nil {
			panic(err)
		}
		path, _ := os.Getwd()
		filename := createUniqueName("attached_files", "asset_", file.Filename, fileType)
		path = filepath.Join(path, "assets/attached_files", filename)
		out, _ := os.Create(path)
		defer out.Close()

		_, err = io.Copy(out, tempFile)

		if err != nil {
			panic(err)
		}
		PRepositoryDB.AttachFiles(pid, filename)
	}

}

// Index is check function.
func Index(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./projectCheck.html"))
	temp.Execute(w, nil)
}

// CheckTemplate is check function.
func CheckTemplate(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseGlob("templates/*.html"))
	temp.ExecuteTemplate(w, "CreateProfile.html", nil)
}

// createUniqueName is a function that creates a unique name for a given asset.
func createUniqueName(tableName string, prefix string, filename string, fileType string) string {

	count := PRepositoryDB.CountMember(tableName)
	suffix := ""
	endPoint := 0

	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			endPoint = i
			break
		}
	}

	for ; endPoint < len(filename); endPoint++ {
		suffix += string(filename[endPoint])
	}

	newFilename := fmt.Sprintf(prefix+fileType+"_%d"+suffix, count)

	for PRepositoryDB.SearchMember(tableName, newFilename) {
		count++
		newFilename = fmt.Sprintf(prefix+fileType+"_%d"+suffix, count)
	}
	return newFilename
}
