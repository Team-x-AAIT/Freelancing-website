package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/client/handler"
	ma "github.com/Team-x-AAIT/Freelancing-website/delivery/http/handler"

	_ "github.com/go-sql-driver/mysql"
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
	var funcMap = template.FuncMap{"ToWorkType": ma.ChangeToWorkType, "ToSDate": ma.ChangeToStandardDate}
	handler.Temp = template.New("").Funcs(funcMap)
	handler.Temp, _ = handler.Temp.ParseGlob("ui/templates/*.html")

	// handler.Temp = template.Must(template.ParseGlob("ui/templates/*.html"))

}

func main() {

	mux := http.NewServeMux()

	fileServer1 := http.FileServer(http.Dir("ui/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer1))

	fileServer := http.FileServer(http.Dir("ui/templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fileServer))

	mux.HandleFunc("/", handler.Dashboard)

	mux.HandleFunc("/ViewProject", handler.ViewProject)
	mux.HandleFunc("/PostProject", handler.PostProject)
	mux.HandleFunc("/UpdateProject", handler.UpdateProject)
	mux.HandleFunc("/RemoveProject", handler.RemoveProject)
	mux.HandleFunc("/UserProjects", handler.GetUserProjects)
	mux.HandleFunc("/user/project/search", handler.SearchProject)
	mux.HandleFunc("/ViewProjectApplicants", handler.ViewApplicationRequests)
	mux.HandleFunc("/ViewApplication", handler.ViewApplication)
	mux.HandleFunc("/Apply", handler.ApplyForProject)
	mux.HandleFunc("/user/project/hire", handler.HireApplicant)
	mux.HandleFunc("/RemoveApplication", handler.RemoveApplication)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()

}
