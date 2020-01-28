package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/application"

	"github.com/julienschmidt/httprouter"

	"github.com/Team-x-AAIT/Freelancing-website/project"
	"github.com/Team-x-AAIT/Freelancing-website/user"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// ProjectHandler is a struct for the project hanlding functions.
type ProjectHandler struct {
	PService project.IService
	AService application.IService
	UService user.IService
}

// NewProjectHandler is a function that return new Project Handler type for the api.
func NewProjectHandler(service project.IService, appService application.IService, uservice user.IService) *ProjectHandler {
	return &ProjectHandler{PService: service,
		AService: appService,
		UService: uservice}
}

// PostProject is a Handler func that initaite project posting process.
func (ph *ProjectHandler) PostProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	uid := params.ByName("uid")

	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	projectHolder := &entity.Project{}

	err := json.Unmarshal(body, projectHolder)

	user := ph.UService.SearchUser(uid)

	if err != nil || user == nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	pid, errMap := ph.PService.PostProject(projectHolder, uid)

	if len(errMap) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("/v1/user/%s/project/%s", uid, pid)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return

}

// UpdateProject is a Handler func that initaite project updating process.
func (ph *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")
	uid := params.ByName("uid")

	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	projectHolder := &entity.Project{}

	err := json.Unmarshal(body, projectHolder)

	user := ph.UService.SearchUser(uid)

	if err != nil || user == nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	tempProject := ph.PService.SearchProjectByID(pid)
	if tempProject.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if !ph.PService.SearchLink(uid, pid) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	projectHolder.ID = pid
	pid, errMap := ph.PService.UpdateProject(projectHolder)
	if len(errMap) > 0 || pid == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return

}

// RemoveProject is a Handler func that initaite project Removing process.
func (ph *ProjectHandler) RemoveProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	uid := params.ByName("uid")
	pid := params.ByName("pid")

	if !ph.PService.SearchLink(uid, pid) {
		fmt.Println("Inside 1")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := ph.AService.RemoveApplicationInfo("", pid); err != nil {
		fmt.Println("Inside 2")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	projectHolder, err := ph.PService.RemoveProjectInformation(uid, pid)
	if err != nil {
		fmt.Println("Inside 3")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(projectHolder, "", "\t\t")
	if err != nil {
		fmt.Println("Inside 4")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// ViewProject is a Handler func that enables applicant to view posted project.
func (ph *ProjectHandler) ViewProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")
	projectW := ph.PService.SearchProjectByID(pid)

	if projectW.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	owner := ph.UService.SearchUser(ph.UService.GetOwner(projectW.ID))
	projectUserContainer := ph.PService.ViewProject(pid, owner, projectW)

	output, err := json.MarshalIndent(projectUserContainer, "", "\t\t")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

// GetUserProjects is a Handler func that sends all the projects linked to a user.
func (ph *ProjectHandler) GetUserProjects(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	uid := params.ByName("uid")
	listOfSentProjects := ph.PService.GetSentProjects(uid)
	output, err := json.MarshalIndent(listOfSentProjects, "", "\t\t")
	if err != nil || len(listOfSentProjects) == 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// SearchProject is a Handler func that initaite projet searching process.
func (ph *ProjectHandler) SearchProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)

	search := &entity.SearchBag{}

	err := json.Unmarshal(body, search)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	filterType, err := strconv.ParseInt(search.FilterTypeS, 0, 0)
	if err != nil {
		filterType = -1
	}
	filterValue1, err := strconv.ParseFloat(search.FilterValue1S, 0)
	if err != nil {
		filterType = 0
	}
	filterValue2, err := strconv.ParseFloat(search.FilterValue2S, 0)
	if err != nil {
		filterType = 0
	}
	pageNum, err := strconv.ParseInt(search.PageNumS, 0, 0)
	if err != nil {
		filterType = 1
	}

	projects := ph.PService.FindProject(search.SearchKey, search.SearchBy, filterType, filterValue1, filterValue2, pageNum)

	output, err := json.MarshalIndent(projects, "", "\t\t")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

// ViewApplicationRequests is a Handler func that enables user to view all the applicants to a certain project.
func (ph *ProjectHandler) ViewApplicationRequests(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")

	listOfApplicantsID := ph.AService.GetProjectApplicantsID(pid)
	var listOfProjectApplicants []entity.Applicants

	for _, value := range listOfApplicantsID {
		applicantInfo := ph.UService.SearchUser(value.ApplicantID)
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
			CreatedAt:    value.CreatedAt,
			Hired:        value.Hired}

		listOfProjectApplicants = append(listOfProjectApplicants, projectApplicant)
	}

	output, err := json.MarshalIndent(listOfProjectApplicants, "", "\t\t")
	if err != nil || len(listOfProjectApplicants) == 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

// ViewApplication is a Handler func that shows an application of a user for a certain user.
func (ph *ProjectHandler) ViewApplication(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")
	applicantUID := params.ByName("uid")

	application, err := ph.AService.GetApplication(applicantUID, pid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(application, "", "\t\t")
	if err != nil {
		fmt.Println("inside + 2")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// ApplyForProject is a Handler func that initaite user application process.
func (ph *ProjectHandler) ApplyForProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")

	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	proposlHolder := &entity.ApplicationBag{}

	err := json.Unmarshal(body, proposlHolder)

	project := ph.PService.SearchProjectByID(pid)

	if ph.PService.SearchLink(proposlHolder.ApplicantID, pid) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil || project.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := ph.AService.Apply(pid, proposlHolder.ApplicantID, proposlHolder.Proposal); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("/v1/user/projects/%s/applications/%s", pid, proposlHolder.ApplicantID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return

}

// HireApplicant is a Handler funct that initaite client hiring processe.
func (ph *ProjectHandler) HireApplicant(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")
	applicantUID := params.ByName("applicant_uid")

	err := ph.AService.HireApplicant(pid, applicantUID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

// RemoveApplication is a Handler func that initaite application removal process
func (ph *ProjectHandler) RemoveApplication(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	pid := params.ByName("pid")
	applicantUID := params.ByName("uid")

	application, err := ph.AService.RemoveApplicant(applicantUID, pid)
	output, err := json.MarshalIndent(application, "", "\t\t")
	if err != nil || application.ApplicantID != applicantUID {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}
