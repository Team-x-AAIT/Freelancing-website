package handler

import (
	"html/template"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/admin"
	"github.com/Team-x-AAIT/Freelancing-website/application"
	apsr "github.com/Team-x-AAIT/Freelancing-website/application"
	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/project"
	"github.com/Team-x-AAIT/Freelancing-website/user"
)

// AdminHandler is a struct for the Admin hanlding functions.
type AdminHandler struct {
	ADService admin.IService
	UService  user.IService
	PService  project.IService
	AService  apsr.IService
	UHandler  *UserHandler
	Temp      *template.Template
}

// NewAdminHandler is a function that return new Admin Handler type.
func NewAdminHandler(service admin.IService, urService user.IService, prService project.IService,
	apService application.IService, uh *UserHandler, temp *template.Template) *AdminHandler {

	return &AdminHandler{ADService: service,
		UService: urService,
		PService: prService,
		AService: apService,
		UHandler: uh,
		Temp:     temp}
}

// RemoveUser is a handler func that removes the user from the system.
func (adh *AdminHandler) RemoveUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		adh.Temp.ExecuteTemplate(w, "remove_user.layout", nil)
	} else if r.Method == http.MethodPost {

		rUID := r.FormValue("rUID")

		userHolder := adh.UService.SearchUser(rUID)
		if userHolder.UID == "" {
			return
		}
		listOfProjects := adh.PService.GetSentProjects("ruid")
		for _, project := range listOfProjects {
			_, err := adh.PService.RemoveProjectInformation(rUID, project.ID)
			if err != nil {
				panic(err)
			}
		}

		listOfApplications := adh.AService.GetUserApplicationHistory(rUID)
		for _, application := range listOfApplications {

			err := adh.AService.RemoveApplicationInfo(rUID, application.PID)
			if err != nil {
				panic(err)
			}
		}

		userHolder, err := adh.UService.RemoveUser(rUID)
		if err != nil {
			panic(err)
		}

		adh.UHandler.Logout(w, r)
		return
	}
}

// SignUp handles signup requests
func (adh *AdminHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		admin := entity.Admin{}
		admin.Firstname = r.FormValue("fname")
		admin.Lastname = r.FormValue("lname")

		admin.Phonenumber = r.FormValue("phone")
		admin.Email = r.FormValue("email")
		admin.Password = r.FormValue("pass")
		_, err := adh.ADService.StoreAdmin(&admin)
		if err != nil {
			adh.Temp.ExecuteTemplate(w, "register.layout", nil)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		adh.Temp.ExecuteTemplate(w, "register.layout", nil)
	}
}

// Index is checking.
func (adh *AdminHandler) Index(w http.ResponseWriter, r *http.Request) {
	adh.Temp.ExecuteTemplate(w, "index.layout", nil)
}

//Update handles Update requests
func (adh *AdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		admin := entity.Admin{}
		admin.AID = r.FormValue("aid")
		admin.Firstname = r.FormValue("fname")
		admin.Lastname = r.FormValue("lname")

		admin.Phonenumber = r.FormValue("phone")
		admin.Email = r.FormValue("email")
		admin.Password = r.FormValue("pass")
		_, err := adh.ADService.UpdateAdmin(&admin)
		if err != nil {
			adh.Temp.ExecuteTemplate(w, "update.layout", nil)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		adh.Temp.ExecuteTemplate(w, "update.layout", nil)
	}
}

//Delete handles delete requests
func (adh *AdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		aid := r.FormValue("aid")
		_, err := adh.ADService.DeleteAdmin(aid)
		if err != nil {
			adh.Temp.ExecuteTemplate(w, "login.layout", nil)
		}
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
	} else {
		adh.Temp.ExecuteTemplate(w, "delete.layout", nil)
	}
}

//RemoveProject handles /removeproject requests
func (adh *AdminHandler) RemoveProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else {
		adh.Temp.ExecuteTemplate(w, "remove_project.layout", nil)
	}
}

//AllUsers handles /allusers requests
func (adh *AdminHandler) AllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else {
		adh.Temp.ExecuteTemplate(w, "all_users.layout", nil)
	}
}

//AllProjects handles /allprojects requests
func (adh *AdminHandler) AllProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else {
		adh.Temp.ExecuteTemplate(w, "all_projects.layout", nil)
	}
}
