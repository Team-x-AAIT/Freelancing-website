package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/Team-x-AAIT/Freelancing-website/api/job"
	"github.com/julienschmidt/httprouter"
)

// declaring a struct
type jobresponse struct {
	Status  string
	Content interface{}
}

// declaring a handler struct
type JobHandler struct {
	jobService job.JobService
}

// init handler
func NewJobHandler(js job.JobService) *JobHandler {
	return &JobHandler{jobService: js}
}

//  api implementation for get jobs
func (jh *JobHandler) GetJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	// email := r.PostFormValue("email")
	search := r.PostFormValue("category")
	// search := r.URL.Query().Get("search")
	fmt.Println(search)

	searchedJob := entity.Job{Category: search}
	job, errs := jh.jobService.Job(&searchedJob)

	if len(errs) > 0 {
		data, err := json.MarshalIndent(&response{Status: "error", Content: nil}, "", "\t")
		if err != nil {

		}
		http.Error(w, string(data), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(response{Status: "success", Content: &job}, "", "\n")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write(output)
	return
}

// api implementation for post job
func (jh *JobHandler) PostJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("inside post job")
	w.Header().Set("Content-type", "application/json")

	job := entity.Job{}

	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)

	err := json.Unmarshal(body, &job)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	j, errs := jh.jobService.StoreJob(&job)

	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p := fmt.Sprintf("/v1/job/%d", j.ID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}

// api implementaion for get jobs
func (jh *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	// email := r.PostFormValue("email")
	// password := r.PostFormValue("password")
	search := r.URL.Query().Get("category")

	// searchedJob := entity.Job{Category: search}
	job, errs := jh.jobService.Jobs(search)
	fmt.Println(job)
	if len(errs) > 0 {
		data, err := json.MarshalIndent(&response{Status: "error", Content: nil}, "", "\t")
		if err != nil {

		}
		http.Error(w, string(data), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(jobresponse{Status: "success", Content: job}, "", "\n")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write(output)
	return
}

// api implementation for get job by id
func (jh *JobHandler) GetJobBy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	// email := r.PostFormValue("email")
	// password := r.PostFormValue("password")
	jo := r.URL.Query().Get("id")
	joint, _ := strconv.ParseInt(jo, 10, 0)
	joid := uint(joint)
	// searchedJob := entity.Job{Category: search}
	job, errs := jh.jobService.JobByID(joid)
	fmt.Println(job)
	if len(errs) > 0 {
		data, err := json.MarshalIndent(&response{Status: "error", Content: nil}, "", "\t")
		if err != nil {

		}
		http.Error(w, string(data), http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(jobresponse{Status: "success", Content: job}, "", "\n")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write(output)
	return
}
