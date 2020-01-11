package main

import (
	"net/http"

	// "github.com/Team-x-AAIT/Freelancing-website/api/entity"

	applyRepository "github.com/Team-x-AAIT/Freelancing-website/api/apply/repository"
	applyService "github.com/Team-x-AAIT/Freelancing-website/api/apply/service"
	"github.com/Team-x-AAIT/Freelancing-website/api/delivery/http/handler"
	jobRepository "github.com/Team-x-AAIT/Freelancing-website/api/job/repository"
	jobService "github.com/Team-x-AAIT/Freelancing-website/api/job/service"
	myjobRepository "github.com/Team-x-AAIT/Freelancing-website/api/myjob/repository"
	myjobService "github.com/Team-x-AAIT/Freelancing-website/api/myjob/service"
	"github.com/Team-x-AAIT/Freelancing-website/api/user/repository"
	"github.com/Team-x-AAIT/Freelancing-website/api/user/service"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func main() {
	dbConn, err := gorm.Open("postgres", "postgres://postgres:admin123@localhost/fjobsdb?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbConn.Close()
	// errs := dbConn.CreateTable(&entity.User{}).GetErrors()

	// if len(errs) > 0 {
	// 	panic(errs)
	// }
	// errs1 := dbConn.CreateTable(&entity.Job{}).GetErrors()

	// if len(errs1) > 0 {
	// 	panic(errs1)
	// }
	// errs2 := dbConn.CreateTable(&entity.MyJob{}).GetErrors()

	// if len(errs2) > 0 {
	// 	panic(errs1)
	// }
	// errs3 := dbConn.CreateTable(&entity.Apply{}).GetErrors()

	// if len(errs3) > 0 {
	// 	panic(errs1)
	// }
	userRepo := repository.NewUserGormRepo(dbConn)
	userServ := service.NewUserService(userRepo)
	userHandl := handler.NewUserHandler(userServ)

	jobRepo := jobRepository.NewJobGormRepo(dbConn)
	jobServ := jobService.NewJobService(jobRepo)
	jobHandl := handler.NewJobHandler(jobServ)

	myjobRepo := myjobRepository.NewMyJobGormRepo(dbConn)
	myjobServ := myjobService.NewMyJobService(myjobRepo)
	myjobHandl := handler.NewMyJobHandler(myjobServ)

	applyRepo := applyRepository.NewApplyGormRepo(dbConn)
	applyServ := applyService.NewApplyService(applyRepo)
	applyHandl := handler.NewApplyHandler(applyServ)
	router := httprouter.New()
	// router.GET("/v1/users", userHandl.GetUsers)
	// router.GET("/v1/users/:id", userHandl.GetSingleUser)
	router.POST("/v1/user", userHandl.GetUser)
	// router.PATCH("/v1/users/:id", userHandl.PatchUser)
	router.POST("/v1/users", userHandl.PostUser)
	router.DELETE("/v1/user/:id", userHandl.DeleteUser)
	router.GET("/v1/user/:id", userHandl.RecommendedJobs)
	// this is for jobs this
	router.GET("/v1/job", jobHandl.GetJob)
	router.GET("/v1/jobs", jobHandl.GetJobs)
	router.POST("/v1/job", jobHandl.PostJob)
	router.GET("/v1/jbyid", jobHandl.GetJobBy)
	// this is for my job
	router.POST("/v1/myjob/:userid/:myjob", myjobHandl.PostMyJob)
	router.GET("/v1/myjob", myjobHandl.GetMyJob)
	// this is for applying for job
	router.POST("/v1/applies", applyHandl.PostApply)
	http.ListenAndServe(":8181", router)
}
