package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/api/apply"
	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
	"github.com/julienschmidt/httprouter"
)

// struct which implements ApplyService
type ApplyHandler struct {
	applyService apply.ApplyService
}

// init apply handler
func NewApplyHandler(as apply.ApplyService) *ApplyHandler {
	return &ApplyHandler{applyService: as}
}

// handler for posting apply all the api implementations goes here
func (ah *ApplyHandler) PostApply(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")

	apply := entity.Apply{}

	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)

	err := json.Unmarshal(body, &apply)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	a, errs := ah.applyService.StoreApply(&apply)

	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	p := fmt.Sprintf("/v1/applies/%d", a.JobID)
	w.Header().Set("Location", p)
	w.WriteHeader(http.StatusCreated)
	return
}
