package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Team-x-AAIT/Freelancing-website/client/service"
)
// handles http requests 
type ApplyHandler struct {
	templ *template.Template
}
// init apply handler
func NewApplyHandler(tmlp *template.Template) *ApplyHandler {
	return &ApplyHandler{templ: tmlp}
}
// hanle apply request
func (ah *ApplyHandler) Apply(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		jo := r.URL.Query().Get("id")
		joint, _ := strconv.ParseInt(jo, 10, 0)
		joid := uint(joint)
		data, err := service.GetJobTo(joid)
		fmt.Println(data)
		if err != nil {
			panic(err)
		}
		ah.templ.ExecuteTemplate(w, "apply.layout", data)
		// http.Redirect(w, r, "http://localhost:8282/indexmain", http.StatusSeeOther)

	}
}
