package main

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Team-x-AAIT/Freelancing-website/api/delivery/http/handler"
	prrp "github.com/Team-x-AAIT/Freelancing-website/project/repository"
	prsr "github.com/Team-x-AAIT/Freelancing-website/project/service"

	urrp "github.com/Team-x-AAIT/Freelancing-website/user/repository"
	ursr "github.com/Team-x-AAIT/Freelancing-website/user/service"

	aprp "github.com/Team-x-AAIT/Freelancing-website/application/repository"
	apsr "github.com/Team-x-AAIT/Freelancing-website/application/service"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	ph *handler.ProjectHandler
)

func init() {

	db, err := sql.Open("mysql", "root:0911@tcp(localhost:3306)/FjobsDB?parseTime=true")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	URepositoryDB := urrp.NewUserRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepositoryDB := prrp.NewProjectRepository(db)
	PService := prsr.NewProjectService(PRepositoryDB)

	ARepositoryDB := aprp.NewApplicationRepository(db)
	AService := apsr.NewApplicationService(ARepositoryDB)

	ph = handler.NewProjectHandler(PService, AService, UService)

}

func main() {

	router := httprouter.New()

	// Route for searching projects using search keys.
	// router.GET("/v1/user/projects/search", ph.SearchProject)
	// Route for viewing a particular project.
	router.GET("/v1/user/:uid/projects/:pid", ph.ViewProject)
	// Route for a getting all the projects created by a particular user.
	router.GET("/v1/user/:uid/projects", ph.GetUserProjects)
	// Route for creating a project by a particular user.
	router.POST("/v1/user/:uid/projects", ph.PostProject)
	// Route for updating a project created by a particular user.
	router.PUT("/v1/user/:uid/projects/:pid", ph.UpdateProject)
	// Route for deleting a particular project.
	router.DELETE("/v1/user/:uid/projects/:pid", ph.RemoveProject)

	router.GET("/v1/user/:uid/projects/:pid/applications", ph.ViewApplicationRequests)
	router.GET("/v1/user/:uid/projects/:pid/application", ph.ViewApplication)
	router.POST("/v1/user/:uid/projects/:pid/applications", ph.ApplyForProject)
	// router.PUT("/v1/user/projects/:pid/applications/:uid", ph.HireApplicant)
	router.DELETE("/v1/user/:uid/projects/:pid/applications", ph.RemoveApplication)

	server := http.Server{
		Addr:    ":8181",
		Handler: router,
	}
	server.ListenAndServe()

}

// Search_Project?searchKey=mobile developer&searchBy=te&filterType=1&filterValue1=100&filterValue2=500&pageNum=1
