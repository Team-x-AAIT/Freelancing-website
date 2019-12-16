package main

import (
	"database/sql"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/user"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {

	db, err := sql.Open("mysql", "root:0911@tcp(localhost:3306)/FjobsDB")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	user.URepositoryDB = user.NewRepository(db)
	user.UService = user.NewService(user.URepositoryDB)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", user.HandleMain)
	mux.HandleFunc("/Register", user.Register)
	mux.HandleFunc("/Verify", user.Verify)
	mux.HandleFunc("/Check_Your_Email", user.Index)

	mux.HandleFunc("/GoogleLogin", user.HandleGoogleLogin)
	mux.HandleFunc("/GoogleCallback", user.HandleGoogleCallback)
	mux.HandleFunc("/LinkedInLogin", user.HandleLinkedInLogin)
	mux.HandleFunc("/LinkedInCallback", user.HandleLinkedInCallback)
	mux.HandleFunc("/FacebookLogin", user.HandleFacebookLogin)
	mux.HandleFunc("/FacebookCallback", user.HandleFacebookCallback)
	mux.HandleFunc("/check", user.Index)
	server := http.Server{
		Addr:    ":1234",
		Handler: mux,
	}
	server.ListenAndServe()
}
