package main

import (
	"database/sql"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/project"
	"github.com/Team-x-AAIT/Freelancing-website/user"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {

	db, err := sql.Open("mysql", "root:0911@tcp(localhost:3306)/FjobsDB?parseTime=true")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	user.URepositoryDB = user.NewRepository(db)
	user.UService = user.NewService(user.URepositoryDB)

	project.PRepositoryDB = project.NewRepository(db)
	project.PService = project.NewService(project.PRepositoryDB)

}

func main() {

	mux := http.NewServeMux()

	fileServer1 := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer1))

	fileServer := http.FileServer(http.Dir("templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fileServer))

	mux.HandleFunc("/", user.HandleMain)
	mux.HandleFunc("/Register", user.Register)
	mux.HandleFunc("/Login", user.Login)
	mux.HandleFunc("/EditProfile/Update", user.UpdateProfile)
	mux.HandleFunc("/Verify", user.Verify)
	mux.HandleFunc("/Dashboard", user.Dashboard)

	mux.HandleFunc("/Login_Page", user.IndexNot)
	mux.HandleFunc("/Logout", user.Logout)

	mux.HandleFunc("/Post_Project", project.PostProject)
	mux.HandleFunc("/Create_Project", project.OnCreateProject)
	mux.HandleFunc("/Review_Project", project.ReviewProject)
	mux.HandleFunc("/asset_viewer", project.ViewAsset)
	mux.HandleFunc("/Update_Project", project.UpdateProject)
	mux.HandleFunc("/Remove_Project", project.RemoveProject)
	mux.HandleFunc("/View_Project", project.ViewProject)
	mux.HandleFunc("/ApplyFor_Project", project.ApplyForProject)
	mux.HandleFunc("/Search_Project", project.SearchProject)

	mux.HandleFunc("/View_Applicants", project.ViewApplicationRequest)
	mux.HandleFunc("/Hire_Applicant", project.HireApplicant)
	mux.HandleFunc("/Remove_Application", project.RemoveApplication)

	mux.HandleFunc("/EditProfile", user.EditProfile)
	mux.HandleFunc("/Check_Your_Email", user.Index)

	mux.HandleFunc("/Add_Match_Tag", user.AddMatchTag)
	mux.HandleFunc("/Remove_Match_Tag", user.RemoveMatchTag)
	mux.HandleFunc("/Get_Match_Tag", user.GetMatchTags)
	mux.HandleFunc("/Get_Projects", user.GetProjectsWMatchTags)
	mux.HandleFunc("/Get_Sent_Projects", project.GetSentProjects)

	mux.HandleFunc("/GoogleLogin", user.HandleGoogleLogin)
	mux.HandleFunc("/GoogleCallback", user.HandleGoogleCallback)
	mux.HandleFunc("/LinkedInLogin", user.HandleLinkedInLogin)
	mux.HandleFunc("/LinkedInCallback", user.HandleLinkedInCallback)
	mux.HandleFunc("/FacebookLogin", user.HandleFacebookLogin)
	mux.HandleFunc("/FacebookCallback", user.HandleFacebookCallback)
	mux.HandleFunc("/check", user.Index)

	mux.HandleFunc("/CheckTemplate", project.CheckTemplate)

	server := http.Server{
		Addr:    ":1234",
		Handler: mux,
	}
	server.ListenAndServe()

}

// Search_Project?searchKey=mobile developer&searchBy=te&filterType=1&filterValue1=100&filterValue2=500&pageNum=1
