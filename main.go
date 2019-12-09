package main

import (
	"database/sql"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
	"github.com/Team-x-AAIT/Freelancing-website/user"
	_ "github.com/go-sql-driver/mysql"
)

var db, err = sql.Open("mysql", "root:0911@tcp(localhost:7777)/FjobsDB")

var repositoryDB = user.NewPsqlUserRepository(db)
var service = user.NewUserService(repositoryDB)

func register(w http.ResponseWriter, r *http.Request) {

	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	password := r.FormValue("password")
	phonenumber := r.FormValue("phonenumber")
	email := r.FormValue("email")
	jobTitle := r.FormValue("jobTitle")
	country := r.FormValue("country")
	city := r.FormValue("city")
	gender := r.FormValue("gender")

	user := entities.NewUser(firstname, lastname, password, phonenumber, email, jobTitle, country, city, gender)

	err := service.RegisterUser(user)
	if err != nil {
		panic(err)
	}
}

func main() {

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/Register", register)
	server := http.Server{
		Addr:    ":1234",
		Handler: mux,
	}
	server.ListenAndServe()
}
