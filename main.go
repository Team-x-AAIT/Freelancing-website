package main

import (
	"database/sql"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
	"github.com/Team-x-AAIT/Freelancing-website/user"
	_ "github.com/go-sql-driver/mysql"
)

var db, err = sql.Open("mysql", "root:0911@tcp(localhost:3336)/FjobsDB")

var repositoryDB = user.NewPsqlUserRepository(db)
var service = user.NewUserService(repositoryDB)

func register(w http.ResponseWriter, r *http.Request) {

	firstname := r.URL.Query().Get("firstname")
	lastname := r.URL.Query().Get("lastname")
	password := r.URL.Query().Get("password")
	phonenumber := r.URL.Query().Get("phonenumber")
	email := r.URL.Query().Get("email")
	jobTitle := r.URL.Query().Get("jobTitle")
	country := r.URL.Query().Get("country")
	city := r.URL.Query().Get("city")
	gender := r.URL.Query().Get("gender")

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
