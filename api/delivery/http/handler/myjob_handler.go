package handler

import (
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/api/myjob"
	"github.com/julienschmidt/httprouter"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
)

// struct to hold myjob response
type myjobresponse struct {
	Status  string
	Content interface{}
}

// struct to implement jobservice
type MyJobHandler struct {
	myjobService myjob.MyJobService
}

// init job handler
func NewMyJobHandler(mjs myjob.MyJobService) *MyJobHandler {
	return &MyJobHandler{myjobService: mjs}
}

// api implementaion for postmyjob on signup
func (mjh *MyJobHandler) PostMyJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// w.Header().Set("Content-type", "application/json")
	fmt.Println("inside post myjob")
	usrid, err := strconv.Atoi(ps.ByName("userid"))
	fmt.Println(usrid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	myusrid := uint(usrid)
	fmt.Println(myusrid)
	job := ps.ByName("myjob")
	fmt.Println(job)
	myjob := entity.MyJob{}
	myjob.Job = job
	myjob.UserID = myusrid

	mj, errs := mjh.myjobService.StoreMyJob(&myjob)

	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	p := fmt.Sprintf("/v1/myjob/%d/%s", mj.UserID, mj.Job)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}

// api implementaion for get my job
func (mjh *MyJobHandler) GetMyJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	usrid, err := strconv.Atoi(r.URL.Query().Get("userid"))
	fmt.Println(usrid)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	job, errs := mjh.myjobService.GetMyJob(usrid)
	fmt.Println(job)
	if len(errs) > 0 {
		data, err := json.MarshalIndent(&myjobresponse{Status: "error", Content: nil}, "", "\t")
		if err != nil {

		}
		http.Error(w, string(data), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(myjobresponse{Status: "success", Content: job}, "", "\t")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write(output)
	return

}
