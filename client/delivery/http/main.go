package main

import (
	"html/template"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/client/delivery/http/handler"
	"github.com/gorilla/mux"
)

var templ = template.Must(template.ParseGlob("../../ui/templates/*"))

func main() {
	userHandler := handler.NewUserHandler(templ)
	applyHandler := handler.NewApplyHandler(templ)
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../../ui/assets"))))
	router.HandleFunc("/", userHandler.Index)
	router.HandleFunc("/indexmain", userHandler.IndexAuth)
	router.HandleFunc("/login", userHandler.Login)
	router.HandleFunc("/signup", userHandler.SignUp)
	router.HandleFunc("/post", userHandler.Post)
	// this is for search
	router.HandleFunc("/search", userHandler.Search)
	router.HandleFunc("/searchauth", userHandler.SearchAuth)
	// this for apply
	router.HandleFunc("/apply", applyHandler.Apply)
	http.ListenAndServe(":8282", router)
}
