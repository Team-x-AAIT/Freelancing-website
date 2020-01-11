package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/Team-x-AAIT/Freelancing-website/client/service"
)
// struct implements useful things for handling requests
type UserHandler struct {
	templ *template.Template
}
// init user handler
func NewUserHandler(tmlp *template.Template) *UserHandler {
	return &UserHandler{templ: tmlp}
}
// handle / get request
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {

	uh.templ.ExecuteTemplate(w, "index.layout", nil)
	// uh.templ.ExecuteTemplate(w, "indexmainauth.layout", nil)
}
// handles users with auth
func (uh *UserHandler) IndexAuth(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("user")
	id, _ := strconv.ParseInt(c.Value, 10, 0)
	// fmt.Println(c.Value)
	// fmt.Println(err)
	myid := uint(id)
	data, err := service.GetRecommendedJobs(myid)
	if err != nil {
		panic(err)
	}
	// uh.templ.ExecuteTemplate(w, "index.layout", nil)
	uh.templ.ExecuteTemplate(w, "indexmainauth.layout", data)
}
// handles sign up request
func (uh *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		uh.templ.ExecuteTemplate(w, "signup.layout", nil)
	} else if r.Method == http.MethodPost {

		firstname := r.PostFormValue("firstname")
		lastname := r.PostFormValue("lastname")
		username := r.PostFormValue("username")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		aboutyou := r.PostFormValue("aboutyou")
		jobs := r.PostFormValue("jobs")

		country := r.PostFormValue("country")

		user := entity.User{FirstName: firstname, LastName: lastname, UserName: username, Email: email, Password: password, AboutYou: aboutyou, Country: country}
		err := service.PostUser(&user)
		if err != nil {
		
			w.Write([]byte("failed"))
			return
		}
		if err == nil {
			// email := r.PostFormValue("email")
			// password := r.PostFormValue("password")
			usr := entity.User{Email: email, Password: password}

			resp, err := service.GetUser(&usr)
			if err != nil {
				panic(err)
			}
			fmt.Println(usr)

			cookie := http.Cookie{
				Name:     "user",
				Value:    strconv.Itoa(int(resp.ID)),
				MaxAge:   60 * 3,
				Path:     "/",
				HttpOnly: true,
			}

			http.SetCookie(w, &cookie)

			myjob := strings.Split(jobs, ",")
			fmt.Println(myjob)
			fmt.Println(usr.ID)
			for _, job := range myjob {
				err := service.PostMyJob(resp.ID, job)
				if err != nil {
					panic(err)
				}
				fmt.Println(job)
			}

			// w.Write([]byte("success"))
			// w.Header().Set("Location:", "http://locahost:8282")
			data, err := service.GetRecommendedJobs(resp.ID)
			fmt.Println(data)
			if err != nil {
				panic(err)
			}
			uh.templ.ExecuteTemplate(w, "indexmainauth.layout", data)
			// http.Redirect(w, r, "http://localhost:8282/indexmain", http.StatusSeeOther)
		}

	}
}
// handles login request
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside Login method..")
	if r.Method == http.MethodGet {
		fmt.Println("inside Login Get method..")
		uh.templ.ExecuteTemplate(w, "login.layout", nil)
	} else if r.Method == http.MethodPost {
		fmt.Println("inside Login Post method..")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		user := entity.User{Email: email, Password: password}
		fmt.Println(user)
		resp, err := service.GetUser(&user)
		if err != nil {
			if err.Error() == "error" {
				uh.templ.ExecuteTemplate(w, "login.layout", "either username or password incorrect")
				return
			}
		} else {

			cookie := http.Cookie{
				Name:     "user",
				Value:    strconv.Itoa(int(resp.ID)),
				MaxAge:   60 * 3,
				Path:     "/",
				HttpOnly: true,
			}

			http.SetCookie(w, &cookie)
			data, err := service.GetRecommendedJobs(resp.ID)
			fmt.Println(data)
			if err != nil {
				panic(err)
			}
			uh.templ.ExecuteTemplate(w, "indexmainauth.layout", data)
			// uh.templ.ExecuteTemplate(w, "indexmainauth.layout", data)

		}
	}
}
// handles search requests
func (uh *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside Search method..")
	if r.Method == http.MethodPost {
		fmt.Println("inside Search Get method..")
		uh.templ.ExecuteTemplate(w, "login.layout", nil)
	} else if r.Method == http.MethodGet {
		fmt.Println("inside Search post method..")
		search := r.FormValue("category")
		fmt.Println(search)
		// jobs := entity.Job{Category: search}
		// fmt.Println(jobs)
		resp, err := service.GetJobs(search)
		if err != nil {
			if err.Error() == "error" {
				uh.templ.ExecuteTemplate(w, "login.layout", "either username or password incorrect")
				return
			}
		} else {
			fmt.Println(resp)
			uh.templ.ExecuteTemplate(w, "indexmainunauth.layout", resp)

		}
	}
}
//handles search for the users with auth
func (uh *UserHandler) SearchAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside Search method..")
	if r.Method == http.MethodPost {
		fmt.Println("inside Search Get method..")
		uh.templ.ExecuteTemplate(w, "login.layout", nil)
	} else if r.Method == http.MethodGet {
		fmt.Println("inside Search post method..")
		search := r.FormValue("category")
		fmt.Println(search)
		// jobs := entity.Job{Category: search}
		// fmt.Println(jobs)
		resp, err := service.GetJobs(search)
		if err != nil {
			if err.Error() == "error" {
				uh.templ.ExecuteTemplate(w, "login.layout", "either username or password incorrect")
				return
			}
		} else {
			fmt.Println(resp)
			uh.templ.ExecuteTemplate(w, "indexmainauth.layout", resp)

		}
	}
}
// handles posting job 
func (uh *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c, err := r.Cookie("user")
		id, _ := strconv.ParseInt(c.Value, 10, 0)
		myid := uint(id)
		data, err := service.GetRecommendedJobs(myid)
		fmt.Println(data)
		if err != nil {
			panic(err)
		}
		uh.templ.ExecuteTemplate(w, "indexmain.layout", data)
	} else if r.Method == http.MethodPost {

		title := r.PostFormValue("title")
		description := r.PostFormValue("description")
		category := r.PostFormValue("category")
		c, err := r.Cookie("user")
		id, _ := strconv.ParseInt(c.Value, 10, 0)
		myid := uint(id)
		job := entity.Job{Title: title, Description: description, Category: category, UserID: myid}
		err = service.PostJob(&job)
		if err != nil {
			w.Write([]byte("failed"))
			return
		}
		if err == nil {

			data, err := service.GetRecommendedJobs(myid)
			fmt.Println(data)
			if err != nil {
				panic(err)
			}
			uh.templ.ExecuteTemplate(w, "indexmainauth.layout", data)
			// http.Redirect(w, r, "http://localhost:8282/indexmain", http.StatusSeeOther)
		}

	}
}
