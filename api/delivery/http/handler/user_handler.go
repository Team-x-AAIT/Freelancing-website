package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/Team-x-AAIT/Freelancing-website/api/user"
	"github.com/julienschmidt/httprouter"
)

// creating response struct
type response struct {
	Status  string
	Content interface{}
}

// creating user handler for implementing userservice
type UserHandler struct {
	userService user.UserService
}

// init user handler
func NewUserHandler(us user.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// get user
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	fmt.Println(email, password)

	usr := entity.User{Email: email, Password: password}
	user, errs := uh.userService.User(&usr)

	if len(errs) > 0 {
		data, err := json.MarshalIndent(&response{Status: "error", Content: nil}, "", "\t")
		if err != nil {

		}
		http.Error(w, string(data), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(response{Status: "success", Content: &user}, "", "\n")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write(output)
	return
}

// delete user
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-type", "application/json")

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	user, errs := uh.userService.DeleteUser(uint(id))

	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = json.MarshalIndent(user, "", "\n")

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return

}

// PostUser requests
func (uh *UserHandler) PostUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")

	user := entity.User{}

	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)

	err := json.Unmarshal(body, &user)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	u, errs := uh.userService.StoreUser(&user)

	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	p := fmt.Sprintf("/v1/users/%d", u.ID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}

func (uh *UserHandler) RecommendedJobs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	fmt.Println("inside recommended..")

	w.Header().Set("Content-type", "application/json")

	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	job, errs := uh.userService.RecommendedJobs(uint(id))

	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(response{Status: "success", Content: job}, "", "\t")

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Write(output)
	w.WriteHeader(http.StatusNoContent)
	fmt.Println(job)
	return

}
