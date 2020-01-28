package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/client/data"
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// Temp contains all the parased templates.
var Temp *template.Template

// InputContainer contains the input value for the templates.
type InputContainer struct {
	ErrorValue  bool
	PUC         *entity.ProjectUserContainer
	Applicants  []*entity.Applicants
	Project     *entity.Project
	Projects    []entity.Project
	Application *entity.ApplicationBag
	Location    string
}

// Dashboard is a Handler func that display the dashboard menu.
func Dashboard(w http.ResponseWriter, r *http.Request) {
	Temp.ExecuteTemplate(w, "ViewProject.html", nil)
}

// ViewProject is a Handler func that enables applicant to view posted project.
func ViewProject(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "ViewProject.html", nil)
		return
	} else if r.Method == http.MethodPost {
		pid := r.FormValue("pid")
		projectUserContainer, err := data.ViewProject(pid)
		input := InputContainer{PUC: projectUserContainer}
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input.ErrorValue = true
			Temp.ExecuteTemplate(w, "ViewProject.html", input)
			return
		}

		Temp.ExecuteTemplate(w, "ViewProject.html", input)
	}

}

// PostProject is a Handler func that initaite project posting process.
func PostProject(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "PostProject.html", nil)
		return
	} else if r.Method == http.MethodPost {
		uid := r.FormValue("uid")
		title := r.FormValue("title")
		description := r.FormValue("description")
		details := r.FormValue("details")
		category := r.FormValue("category")
		subcategory := r.FormValue("subcategory")
		budgetS := r.FormValue("budget")
		worktypeS := r.FormValue("worktype")

		budget, err := strconv.ParseFloat(budgetS, 0)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "PostProject.html", input)
			return
		}
		worktype, err := strconv.ParseInt(worktypeS, 0, 0)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "PostProject.html", input)
			return
		}

		projectHolder := entity.NewProject(title, description, details, category, subcategory, budget, worktype)
		pid, err := data.PostProject(projectHolder, uid)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "PostProject.html", input)
			return
		}
		input := InputContainer{Location: pid}
		Temp.ExecuteTemplate(w, "PostProject.html", input)
	}

}

// UpdateProject is a Handler func that initaite project updating process.
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "UpdateProject.html", nil)
		return
	} else if r.Method == http.MethodPost {
		uid := r.FormValue("uid")
		pid := r.FormValue("pid")
		title := r.FormValue("title")
		description := r.FormValue("description")
		details := r.FormValue("details")
		category := r.FormValue("category")
		subcategory := r.FormValue("subcategory")
		budgetS := r.FormValue("budget")
		worktypeS := r.FormValue("worktype")

		budget, err := strconv.ParseFloat(budgetS, 0)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "UpdateProject.html", input)
			return
		}
		worktype, err := strconv.ParseInt(worktypeS, 0, 0)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "UpdateProject.html", input)
			return
		}

		projectHolder := entity.NewProject(title, description, details, category, subcategory, budget, worktype)

		projectHolder.ID = pid
		err = data.UpdateProject(uid, projectHolder)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "UpdateProject.html", input)
			return
		}

		message := "Update Successful!"
		input := InputContainer{Location: message}
		Temp.ExecuteTemplate(w, "UpdateProject.html", input)
	}

}

// RemoveProject is a Handler func that initaite project Removing process.
func RemoveProject(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "RemoveProject.html", nil)
		return
	} else if r.Method == http.MethodPost {
		uid := r.FormValue("uid")
		pid := r.FormValue("pid")

		project, err := data.RemoveProjectInformation(uid, pid)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "RemoveProject.html", input)
			return
		}
		input := InputContainer{Project: project}
		Temp.ExecuteTemplate(w, "RemoveProject.html", input)
	}
}

// GetUserProjects is a Handler func that sends all the projects linked to a user.
func GetUserProjects(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "UserProjects.html", nil)
		return
	} else if r.Method == http.MethodPost {
		uid := r.FormValue("uid")
		listOfSentProjects, err := data.GetUserProjects(uid)

		if err != nil {
			input := InputContainer{ErrorValue: true}
			// w.WriteHeader(http.StatusNoContent)
			Temp.ExecuteTemplate(w, "UserProjects.html", input)
			return
		}

		input := InputContainer{Projects: listOfSentProjects}
		Temp.ExecuteTemplate(w, "UserProjects.html", input)
	}
}

// SearchProject is a Handler func that initaite projet searching process.
func SearchProject(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "SearchProjects.html", nil)
		return
	} else if r.Method == http.MethodPost {
		searchKey := r.FormValue("searchKey")
		searchBy := r.FormValue("searchBy")
		filterTypeS := r.FormValue("filterType")
		filterValue1S := r.FormValue("filterValue1")
		filterValue2S := r.FormValue("filterValue2")
		pageNumS := r.FormValue("pageNum")

		projects, err := data.FindProject(searchKey, searchBy, filterTypeS, filterValue1S, filterValue2S, pageNumS)

		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			Temp.ExecuteTemplate(w, "error.layout", nil)
			return
		}

		Temp.ExecuteTemplate(w, "SearchProjects.html", projects)
	}

}

// ViewApplicationRequests is a Handler func that enables user to view all the applicants to a certain project.
func ViewApplicationRequests(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "ViewProjectApplicants.html", nil)
		return
	} else if r.Method == http.MethodPost {
		pid := r.FormValue("pid")

		listOfProjectApplicants, err := data.GetProjectApplicants(pid)
		// fmt.Println("-------------------:" + listOfProjectApplicants[0].Email)

		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "ViewProjectApplicants.html", input)
			return
		}

		input := InputContainer{Applicants: listOfProjectApplicants}
		Temp.ExecuteTemplate(w, "ViewProjectApplicants.html", input)
	}
}

// ViewApplication is a Handler func that enables user to view a certain application.
func ViewApplication(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "ViewApplication.html", nil)
		return
	} else if r.Method == http.MethodPost {
		pid := r.FormValue("pid")
		uid := r.FormValue("uid")

		application, err := data.GetApplication(pid, uid)

		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "ViewApplication.html", input)
			return
		}

		input := InputContainer{Application: application}
		Temp.ExecuteTemplate(w, "ViewApplication.html", input)
	}
}

// ApplyForProject is a Handler func that initaite user application process.
func ApplyForProject(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "Apply.html", nil)
		return
	} else if r.Method == http.MethodPost {
		applicantUID := r.FormValue("applicant_uid")
		pid := r.FormValue("pid")
		proposal := r.FormValue("proposal")

		application := new(entity.ApplicationBag)
		application.ApplicantID = applicantUID
		application.PID = pid
		application.Proposal = proposal

		_, err := data.Apply(application)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "Apply.html", input)
			return
		}

		input := InputContainer{Application: application, Location: "new location"}
		Temp.ExecuteTemplate(w, "Apply.html", input)
	}

}

// HireApplicant is a Handler funct that initaite client hiring processe.
func HireApplicant(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "HireApplicant.html", nil)
		return
	} else if r.Method == http.MethodPost {
		pid := r.FormValue("pid")
		applicantUID := r.FormValue("applicant_UID")

		err := data.HireApplicant(pid, applicantUID)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			Temp.ExecuteTemplate(w, "error.layout", nil)
			return
		}

		message := "Applicant hired!"
		Temp.ExecuteTemplate(w, "HireApplicant.html", message)
	}
}

// RemoveApplication is a Handler func that initaite application removal process
func RemoveApplication(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Temp.ExecuteTemplate(w, "RemoveApplication.html", nil)
		return
	} else if r.Method == http.MethodPost {
		pid := r.FormValue("pid")
		applicantUID := r.FormValue("applicant_UID")

		application, err := data.RemoveApplicant(applicantUID, pid)
		if err != nil {
			// w.WriteHeader(http.StatusNoContent)
			input := InputContainer{ErrorValue: true}
			Temp.ExecuteTemplate(w, "RemoveApplication.html", input)
			return
		}

		input := InputContainer{Application: application}
		Temp.ExecuteTemplate(w, "RemoveApplication.html", input)
	}
}
