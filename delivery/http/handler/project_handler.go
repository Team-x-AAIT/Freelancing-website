package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Team-x-AAIT/Freelancing-website/application"

	"github.com/Team-x-AAIT/Freelancing-website/stringTools"

	"github.com/Team-x-AAIT/Freelancing-website/project"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// PInputContainer is a struct that holds information for the templates.
type PInputContainer struct {
	LoggedInUser         *entity.User
	Project              *entity.Project
	Projects             []*entity.Project
	ProjectUserContainer *entity.ProjectUserContainer
	Applications         []*entity.ApplicationBag
	Error                entity.ErrorBag
	Categories           []string
	SubCategories        []string
	Counter              [3]int64
	Prefe                bool
	CSRF                 string
}

// ProjectHandler is a struct for the project hanlding functions.
type ProjectHandler struct {
	PService project.IService
	AService application.IService
	UHandler *UserHandler
	Temp     *template.Template
}

// NewProjectHandler is a function that return new Project Handler type.
func NewProjectHandler(service project.IService, appServices application.IService, uh *UserHandler, temp *template.Template) *ProjectHandler {
	return &ProjectHandler{PService: service,
		AService: appServices,
		UHandler: uh,
		Temp:     temp}
}

// Dashboard is a Handler func that initaite the Home page of a user after checking his/her profile is completed.
func (ph *ProjectHandler) Dashboard(w http.ResponseWriter, r *http.Request) {

	user := ph.UHandler.Authentication(r)
	if user == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	// Means incomplete profile.
	if user.Gender == "" {
		http.Redirect(w, r, "/EditProfile", http.StatusSeeOther)
		return
	}

	inputContainer := PInputContainer{LoggedInUser: user}
	inputContainer.Categories = ph.PService.GetCategories()
	if user.Prefe == 0 {
		inputContainer.Prefe = true
	}

	ph.Temp.ExecuteTemplate(w, "Dashboard.html", inputContainer)

}

// PostProject is a Handler func that initaite project posting process.
func (ph *ProjectHandler) PostProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)

	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	// loggedInUser := &entity.UserMock
	if r.Method == http.MethodGet {

		ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: &entity.Project{}, CSRF: token}
		inputContainer.Categories = ph.PService.GetCategories()
		inputContainer.SubCategories = ph.PService.GetSubCategories()

		ph.Temp.ExecuteTemplate(w, "PostProject.html", inputContainer)
		return

	} else if r.Method == http.MethodPost {

		errMap := make(map[string]string)
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
			errMap["budget"] = "Invalid amount!"
		}
		worktype, err := strconv.ParseInt(worktypeS, 0, 0)
		if err != nil {
			errMap["worktype"] = "Invalid worktype!"
		}

		csrfToken := r.FormValue("csrf")
		ok, errCRFS := stringTools.ValidCSRF(csrfToken, ph.UHandler.CSRF)
		if !ok || errCRFS != nil {
			errMap["csrf"] = "Invalid token used!"
		}

		projectHolder := entity.NewProject(title, description, details, category, subcategory, budget, worktype)

		if len(errMap) > 0 {
			ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: projectHolder, Error: errMap, CSRF: token}
			inputContainer.Categories = ph.PService.GetCategories()
			inputContainer.SubCategories = ph.PService.GetSubCategories()

			ph.Temp.ExecuteTemplate(w, "PostProject.html", inputContainer)
			return
		}

		err = ph.VerifyResource("files", "project", r)

		if err != nil {
			switch err.Error() {
			case "invalid format":
				errMap["cv"] = "Only pdf format is allowed!"
			case "file too large":
				errMap["cv"] = "File size should be less than 5MB!"
			case "too many files":
				errMap["cv"] = "should not upload more than 3 files!"
			case "unable to read data":
				errMap["cv"] = "Unable to read data!"
			}

			ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, _ := stringTools.CSRFToken(ph.UHandler.CSRF)
			inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: projectHolder, Error: errMap, CSRF: token}
			inputContainer.Categories = ph.PService.GetCategories()
			inputContainer.SubCategories = ph.PService.GetSubCategories()

			ph.Temp.ExecuteTemplate(w, "PostProject.html", inputContainer)
			return
		}

		pid, errMap := ph.PService.PostProject(projectHolder, uid)

		if len(errMap) > 0 {
			ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: projectHolder, Error: errMap, CSRF: token}
			inputContainer.Categories = ph.PService.GetCategories()
			inputContainer.SubCategories = ph.PService.GetSubCategories()

			ph.Temp.ExecuteTemplate(w, "PostProject.html", inputContainer)
			return
		}

		ph.ResourceExtractor(pid, "files", "project", r)
		http.Redirect(w, r, "/Dashboard", http.StatusSeeOther)
	}

}

// ReviewProject is a Handler func that display the created project for editing.
func (ph *ProjectHandler) ReviewProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	pid := r.FormValue("pid")

	project := ph.PService.SearchProjectByID(pid)

	ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
	token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: project, CSRF: token}
	inputContainer.Categories = ph.PService.GetCategories()
	inputContainer.SubCategories = ph.PService.GetSubCategories()

	ph.Temp.ExecuteTemplate(w, "ReviewProject.html", inputContainer)

}

// UpdateProject is a Handler func that initaite project updating process.
func (ph *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

	errMap := make(map[string]string)
	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
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
		errMap["budget"] = "Invalid amount!"
	}
	worktype, err := strconv.ParseInt(worktypeS, 0, 0)
	if err != nil {
		errMap["worktype"] = "Invalid worktype!"
	}

	csrfToken := r.FormValue("csrf")
	ok, errCRFS := stringTools.ValidCSRF(csrfToken, ph.UHandler.CSRF)
	if !ok || errCRFS != nil {
		errMap["csrf"] = "Invalid token used!"
	}

	projectHolder := entity.NewProject(title, description, details, category, subcategory, budget, worktype)
	projectHolder.ID = pid

	if len(errMap) > 0 {
		ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: projectHolder, Error: errMap, CSRF: token}
		inputContainer.Categories = ph.PService.GetCategories()
		inputContainer.SubCategories = ph.PService.GetSubCategories()

		ph.Temp.ExecuteTemplate(w, "ReviewProject.html", inputContainer)
		return
	}

	err = ph.VerifyResource("files", "project", r)

	if err != nil {
		switch err.Error() {
		case "invalid format":
			errMap["cv"] = "Only pdf format is allowed!"
		case "file too large":
			errMap["cv"] = "File size should be less than 5MB!"
		case "too many files":
			errMap["cv"] = "should not upload more than 3 files!"
		case "unable to read data":
			errMap["cv"] = "Unable to read data!"
		}

		ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, _ := stringTools.CSRFToken(ph.UHandler.CSRF)
		inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: projectHolder, Error: errMap, CSRF: token}
		inputContainer.Categories = ph.PService.GetCategories()
		inputContainer.SubCategories = ph.PService.GetSubCategories()

		ph.Temp.ExecuteTemplate(w, "ReviewProject.html", inputContainer)
		return
	}

	pid, errMap = ph.PService.UpdateProject(projectHolder)

	if len(errMap) > 0 {
		ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		inputContainer := PInputContainer{LoggedInUser: loggedInUser, Project: projectHolder, Error: errMap, CSRF: token}
		inputContainer.Categories = ph.PService.GetCategories()
		inputContainer.SubCategories = ph.PService.GetSubCategories()

		ph.Temp.ExecuteTemplate(w, "ReviewProject.html", inputContainer)
		return
	}

	ph.ResourceExtractor(pid, "files", "project", r)

	var beginIndex int
	for index, char1 := range attachedFilesR {
		if char1 == ',' {
			filename := attachedFilesR[beginIndex:index]
			if err = ph.RemoveAsset(pid, filename); err != nil {
				panic(err)
			}
			beginIndex = index + 1
		}
	}

	http.Redirect(w, r, "/Dashboard", http.StatusSeeOther)

}

// RemoveProject is a Handler func that initaite project Removing process.
func (ph *ProjectHandler) RemoveProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	uid := loggedInUser.UID
	pid := r.FormValue("pid")

	if err := ph.AService.RemoveApplicationInfo("", pid); err != nil {
		panic(err)
	}

	if _, err := ph.PService.RemoveProjectInformation(uid, pid); err != nil {
		panic(err)
	}

}

// ViewProject is a Handler func that enables applicant to view posted project.
func (ph *ProjectHandler) ViewProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	pid := r.FormValue("pid")
	projectW := ph.PService.SearchProjectByID(pid)
	owner := ph.UHandler.UService.SearchUser(ph.UHandler.UService.GetOwner(projectW.ID))
	projectUserContainer := ph.PService.ViewProject(pid, owner, projectW)
	value := PInputContainer{LoggedInUser: loggedInUser, ProjectUserContainer: projectUserContainer}
	ph.Temp.ExecuteTemplate(w, "ViewProject.html", value)

}

// GetSentProjects is a Handler func that sends all the projects linked to a user.
func (ph *ProjectHandler) GetSentProjects(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	listOfSentProjects := ph.PService.GetSentProjects(loggedInUser.UID)
	json, err := json.Marshal(listOfSentProjects)
	if err != nil {
		panic(err)
	}

	w.Write(json)
}

// SearchProject is a Handler func that initaite projet searching process.
func (ph *ProjectHandler) SearchProject(w http.ResponseWriter, r *http.Request) {

	errMap := make(map[string]string)
	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {

		ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		inputContainer := PInputContainer{LoggedInUser: loggedInUser, CSRF: token}

		ph.Temp.ExecuteTemplate(w, "SearchProjects.html", inputContainer)
		return

	} else if r.Method == http.MethodPost {
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
			filterValue1 = 0
		}
		filterValue2, err := strconv.ParseFloat(filterValue2S, 0)
		if err != nil {
			filterValue2 = 0
		}
		pageNum, err := strconv.ParseInt(pageNumS, 0, 0)
		if err != nil {
			pageNum = 1
		}

		csrfToken := r.FormValue("csrf")
		ok, errCRFS := stringTools.ValidCSRF(csrfToken, ph.UHandler.CSRF)
		if !ok || errCRFS != nil {
			errMap["csrf"] = "Invalid token used!"
		}

		if len(errMap) > 0 {
			ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			inputContainer := PInputContainer{LoggedInUser: loggedInUser, Error: errMap, CSRF: token}

			ph.Temp.ExecuteTemplate(w, "SearchProjects.html", inputContainer)
			return
		}

		projects := ph.PService.FindProject(searchKey, searchBy, filterType, filterValue1, filterValue2, pageNum)

		ph.UHandler.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(ph.UHandler.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		inputContainer := PInputContainer{LoggedInUser: loggedInUser, Projects: projects, Error: errMap, CSRF: token}
		ph.Temp.ExecuteTemplate(w, "SearchProjects.html", inputContainer)
	}

}

// ViewApplicationRequest is a Handler func that enables user to view all the applicants to a certain project.
func (ph *ProjectHandler) ViewApplicationRequest(w http.ResponseWriter, r *http.Request) {

	pid := r.FormValue("pid")
	listOfApplicantsID := ph.AService.GetProjectApplicantsID(pid)
	var listOfProjectApplicants []entity.Applicants

	for _, value := range listOfApplicantsID {
		applicantInfo := ph.UHandler.UService.SearchUser(value.ApplicantID)
		projectApplicant := entity.Applicants{
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

// ViewApplicationsStatus is used for viewing all the application sent by the user.
func (ph *ProjectHandler) ViewApplicationsStatus(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	applications := ph.AService.GetUserApplicationHistory(loggedInUser.UID)

	for _, value := range applications {
		project := ph.PService.SearchProjectByID(value.PID)
		value.Project = project
		value.Status = ph.AService.CheckApplicationStatus(project, value.ApplicantID)

	}

	inputContainer := PInputContainer{LoggedInUser: loggedInUser, Applications: applications}
	ph.Temp.ExecuteTemplate(w, "AppliedFor.html", inputContainer)
}

// ApplyForProject is a Handler func that initaite user application process.
func (ph *ProjectHandler) ApplyForProject(w http.ResponseWriter, r *http.Request) {

	loggedInUser := ph.UHandler.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	uid := loggedInUser.UID
	pid := r.FormValue("pid")
	proposal := r.FormValue("proposal")

	if ph.PService.SearchLink(uid, pid) {
		w.Write([]byte(errors.New("Can't apply for own project").Error()))
	}

	if err := ph.AService.Apply(pid, uid, proposal); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("okay"))
	}

}

// HireApplicant is a Handler funct that initaite client hiring processe.
func (ph *ProjectHandler) HireApplicant(w http.ResponseWriter, r *http.Request) {

	pid := r.FormValue("pid")
	applicantUID := r.FormValue("applicant_UID")

	err := ph.AService.HireApplicant(pid, applicantUID)
	if err != nil {
		panic(err)
	}
}

// RemoveApplication is a Handler func that initaite application removal process
func (ph *ProjectHandler) RemoveApplication(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("pid")
	applicantUID := r.FormValue("applicant_UID")

	_, err := ph.AService.RemoveApplicant(applicantUID, pid)
	if err != nil {
		panic(err)
	}
}

// ViewAsset is used for diplaying assets.
func (ph *ProjectHandler) ViewAsset(w http.ResponseWriter, r *http.Request) {
	asset := r.FormValue("asset")
	folder := r.FormValue("type")
	if folder == "cv" {
		asset = "./assets/cv/" + asset
	}
	if folder == "attached_file" {
		asset = "./assets/attached_files/" + asset
	}

	ph.Temp.ExecuteTemplate(w, "assetViewer.html", asset)
}

// RemoveAsset is a function that removes a given file path from the system.
func (ph *ProjectHandler) RemoveAsset(pid string, filename string) error {

	if err := ph.PService.RemoveAttachedFile(pid, filename); err != nil {
		return err
	}
	if err := ph.PService.RemoveFile(filename); err != nil {
		return err
	}

	return nil

}

// VerifyResource will verify the resource is valid before extraction.
func (ph *ProjectHandler) VerifyResource(name string, fileType string, r *http.Request) error {

	err := r.ParseMultipartForm(200000)
	if err != nil {
		panic(err)
	}
	formdata := r.MultipartForm
	files := formdata.File[name]

	for index, file := range files {
		OFile, err := file.Open()
		if err != nil {
			return errors.New("unable to read data")
		}
		defer OFile.Close()
		tempFile, _ := ioutil.ReadAll(OFile)
		tempFileType := http.DetectContentType(tempFile)

		if tempFileType != "application/pdf" {
			return errors.New("invalid format")
		}

		if file.Size > 5000000 {
			return errors.New("file to large")
		}

		if index > 3 {
			return errors.New("too many files")
		}
	}
	return nil

}

// ResourceExtractor is a function that extract file from a request.
func (ph *ProjectHandler) ResourceExtractor(pid string, name string, fileType string, r *http.Request) {

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
		filename := ph.createUniqueName("attached_files", "asset_", file.Filename, fileType)
		path = filepath.Join(path, "..", "..", "ui", "assets/attached_files", filename)
		out, _ := os.Create(path)
		defer out.Close()

		_, err = io.Copy(out, tempFile)

		if err != nil {
			panic(err)
		}
		ph.PService.AttachFiles(pid, filename)
	}

}

// Index is check function.
func (ph *ProjectHandler) Index(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./projectCheck.html"))
	temp.Execute(w, nil)
}

// CheckTemplate is check function.
func (ph *ProjectHandler) CheckTemplate(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseGlob("templates/*.html"))
	temp.ExecuteTemplate(w, "CreateProfile.html", nil)
}

// GetSubCategories is a Handler func that sends a json file containing all the subcategories related to a category.
func (ph *ProjectHandler) GetSubCategories(w http.ResponseWriter, r *http.Request) {

	category := r.FormValue("category")
	listOfSubCategories := ph.PService.GetSubCategoriesOf(category)
	jsonString, _ := json.Marshal(listOfSubCategories)
	w.Write([]byte(jsonString))

}

// createUniqueName is a function that creates a unique name for a given asset.
func (ph *ProjectHandler) createUniqueName(tableName string, prefix string, filename string, fileType string) string {

	count := ph.PService.CountMember(tableName)
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

	for ph.PService.SearchMember(tableName, newFilename) {
		count++
		newFilename = fmt.Sprintf(prefix+fileType+"_%d"+suffix, count)
	}
	return newFilename
}

// ChangeToWorkType is a function that changes a work id to a worktype string.
func ChangeToWorkType(workID int64) string {

	switch workID {
	case 1:
		return "Fixed"
	case 2:
		return "Perhour"
	case 3:
		return "Negotiable"
	}

	return ""
}

// GetStatus is a function that changes a status id to status string
func GetStatus(status int64) string {

	switch status {
	case 0:
		return "Pending"
	case 1:
		return "Hired"
	case 2:
		return "Rejected"
	case 3:
		return "Fired"
	case 4:
		return "Removed"
	}

	return ""
}

// GetColor is a function that returns the color class of a status.
func GetColor(status int64) string {

	switch status {
	case 0:
		return "yellow"
	case 1:
		return "green"
	case 2:
		return "red"
	case 3:
		return "red"
	case 4:
		return "red"
	}

	return ""
}

// ChangeToStandardDate change a time type to standard string date formate.
func ChangeToStandardDate(t time.Time) string {
	layoutUS := "January 2, 2006"
	return t.Format(layoutUS)
}
