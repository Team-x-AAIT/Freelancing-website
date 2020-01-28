package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/delivery/http/handler"
	prrp "github.com/Team-x-AAIT/Freelancing-website/project/repository"
	prsr "github.com/Team-x-AAIT/Freelancing-website/project/service"

	urrp "github.com/Team-x-AAIT/Freelancing-website/user/repository"
	ursr "github.com/Team-x-AAIT/Freelancing-website/user/service"

	aprp "github.com/Team-x-AAIT/Freelancing-website/application/repository"
	apsr "github.com/Team-x-AAIT/Freelancing-website/application/service"

	adrp "github.com/Team-x-AAIT/Freelancing-website/admin/repository"
	adsr "github.com/Team-x-AAIT/Freelancing-website/admin/service"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	uh *handler.UserHandler
	ph *handler.ProjectHandler
	ah *handler.AdminHandler
)

func init() {

	db, err := sql.Open("mysql", "root:0911@tcp(localhost:3306)/FjobsDB?parseTime=true")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Temp contains all the parsed templates.
	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor, "ToSDate": handler.ChangeToStandardDate}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	Temp2 := template.Must(template.ParseGlob("../../ui/templates/adtemplates/templates/*.html"))

	URepositoryDB := urrp.NewUserRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepositoryDB := prrp.NewProjectRepository(db)
	PService := prsr.NewProjectService(PRepositoryDB)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp2)

}

func main() {

	mux := http.NewServeMux()

	fileServer1 := http.FileServer(http.Dir("../../ui/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer1))

	fileServer := http.FileServer(http.Dir("../../ui/templates/"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fileServer))

	fileServer2 := http.FileServer(http.Dir("../../ui/templates/adtemplates/templates/"))
	mux.Handle("/adtemplates/templates/", http.StripPrefix("/templates/", fileServer2))

	mux.HandleFunc("/", uh.WelcomePage)
	mux.HandleFunc("/Register", uh.Register)
	mux.HandleFunc("/Login", uh.Login)
	mux.HandleFunc("/EditProfile", uh.UpdateProfile)
	mux.HandleFunc("/Verify", uh.Verify)
	mux.HandleFunc("/Dashboard", ph.Dashboard)
	mux.HandleFunc("/Logout", uh.Logout)

	mux.HandleFunc("/Remove_User", ah.RemoveUser)

	mux.HandleFunc("/Post_Project", ph.PostProject)
	mux.HandleFunc("/Review_Project", ph.ReviewProject)
	mux.HandleFunc("/asset_viewer", ph.ViewAsset)
	mux.HandleFunc("/Update_Project", ph.UpdateProject)
	mux.HandleFunc("/Remove_Project", ph.RemoveProject)
	mux.HandleFunc("/View_Project", ph.ViewProject)
	mux.HandleFunc("/ApplyFor_Project", ph.ApplyForProject)
	mux.HandleFunc("/Search_Project", ph.SearchProject)

	mux.HandleFunc("/View_Applications_Status", ph.ViewApplicationsStatus)
	mux.HandleFunc("/View_Applicants", ph.ViewApplicationRequest)
	mux.HandleFunc("/Hire_Applicant", ph.HireApplicant)
	mux.HandleFunc("/Remove_Application", ph.RemoveApplication)
	mux.HandleFunc("/SubCategories", ph.GetSubCategories)

	mux.HandleFunc("/Check_Your_Email", uh.Index)

	mux.HandleFunc("/Add_Match_Tag", uh.AddMatchTag)
	mux.HandleFunc("/Remove_Match_Tag", uh.RemoveMatchTag)
	mux.HandleFunc("/Get_Match_Tag", uh.GetMatchTags)
	mux.HandleFunc("/Get_Projects", uh.GetProjectsWMatchTags)
	mux.HandleFunc("/Get_Sent_Projects", ph.GetSentProjects)

	mux.HandleFunc("/GoogleLogin", uh.HandleGoogleLogin)
	mux.HandleFunc("/GoogleCallback", uh.HandleGoogleCallback)
	mux.HandleFunc("/LinkedInLogin", uh.HandleLinkedInLogin)
	mux.HandleFunc("/LinkedInCallback", uh.HandleLinkedInCallback)
	mux.HandleFunc("/FacebookLogin", uh.HandleFacebookLogin)
	mux.HandleFunc("/FacebookCallback", uh.HandleFacebookCallback)
	mux.HandleFunc("/check", uh.Index)

	// http.HandleFunc("/signup",adminHandler.SignUp)
	// 	http.HandleFunc("/login", adminHandler.LogIn)
	// 	http.HandleFunc("/", adminHandler.Index)
	// 	http.HandleFunc("/update", adminHandler.Update)
	// 	http.HandleFunc("/delete", adminHandler.Delete)
	mux.HandleFunc("/admin/removeuser", ah.RemoveUser)
	// 	http.HandleFunc("/removeproject", adminHandler.RemoveProject)
	// 	http.HandleFunc("/allusers", adminHandler.AllUsers)
	// 	http.HandleFunc("/allprojects", adminHandler.AllProjects)

	mux.HandleFunc("/CheckTemplate", ph.CheckTemplate)

	server := http.Server{
		Addr:    ":1234",
		Handler: mux,
	}
	server.ListenAndServe()

}
