package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Team-x-AAIT/Freelancing-website/delivery/http/handler"
	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/project/repository"
	"github.com/Team-x-AAIT/Freelancing-website/project/service"

	urrp "github.com/Team-x-AAIT/Freelancing-website/user/repository"
	ursr "github.com/Team-x-AAIT/Freelancing-website/user/service"

	aprp "github.com/Team-x-AAIT/Freelancing-website/application/repository"
	apsr "github.com/Team-x-AAIT/Freelancing-website/application/service"

	adrp "github.com/Team-x-AAIT/Freelancing-website/admin/repository"
	adsr "github.com/Team-x-AAIT/Freelancing-website/admin/service"
)

func TestPostProject(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Post_Project", ph.PostProject)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar
	form := url.Values{}
	form.Add("Title", entity.ProjectMock.Title)
	form.Add("Description", entity.ProjectMock.Description)
	form.Add("Details", entity.ProjectMock.Details)
	//form.Add("AttachedFiles", entity.ProjectMock.AttachedFiles)
	form.Add("Category", entity.ProjectMock.Category)
	form.Add("Subcategory", entity.ProjectMock.Subcategory)
	//form.Add("Budget", entity.ProjectMock.Budget)
	//form.Add("WorkType", entity.ProjectMock.WorkType)
	//form.Add("Closed", entity.ProjectMock.Closed)

	resp, err := client.PostForm(urlString+"/Post_Project", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("Post Project")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestReviewProject(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Review_Project", ph.PostProject)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar

	resp, err := client.Get(urlString + "/Review_Project?pid=PID1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	valueToCompare := `Get your project Done!`

	if !bytes.Contains(body, []byte(valueToCompare)) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestUpdateProject(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Update_Project", ph.PostProject)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar

	resp, err := client.Get(urlString + "/Update_Project")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	valueToCompare := `Get your project Done!`

	if !bytes.Contains(body, []byte(valueToCompare)) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestRemoveProject(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Remove_Project", ph.PostProject)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar

	resp, err := client.Get(urlString + "/Remove_Project")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	valueToCompare := `Get your project Done`

	if !bytes.Contains(body, []byte(valueToCompare)) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestDashboard(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Dashboard", ph.Dashboard)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar

	resp, err := client.Get(urlString + "/Dashboard")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	valueToCompare := `Welcome To FJob`

	if !bytes.Contains(body, []byte(valueToCompare)) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestSearchProject(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Search_Project", ph.SearchProject)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar
	form := url.Values{}
	form.Add("search_by", "sfkaljhfasf")
	form.Add("title", "web")
	//form.Add("worktype",1)
	resp, err := client.PostForm(urlString+"/Search_Project", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("Search Project")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestGetSubCategories(t *testing.T) {

	var funcMap = template.FuncMap{"ToWorkType": handler.ChangeToWorkType, "GetStatus": handler.GetStatus, "GetColor": handler.GetColor}
	Temp := template.New("").Funcs(funcMap)
	Temp, _ = Temp.ParseGlob("../../ui/templates/*.html")

	URepositoryDB := urrp.NewUserMockRepository(db)
	UService := ursr.NewUserService(URepositoryDB)

	PRepository := repository.NewMockRepository(nil)
	PService := service.NewProjectService(PRepository)

	ADRepositoryDB := adrp.NewAdminRepository(db)
	ADService := adsr.NewAdminService(ADRepositoryDB)

	APRepositoryDB := aprp.NewApplicationRepository(db)
	APService := apsr.NewApplicationService(APRepositoryDB)

	SRepositoryDB := urrp.NewMockSessionRepository(db)
	SService := ursr.NewSessionService(SRepositoryDB)

	uh = handler.NewUserHandler(UService, SService, Temp, []byte("Protecting_from_CSRF"))
	ph = handler.NewProjectHandler(PService, APService, uh, Temp)
	ah = handler.NewAdminHandler(ADService, UService, PService, APService, uh, Temp)

	ph = handler.NewProjectHandler(PService, APService, uh, Temp)

	mux := http.NewServeMux()
	mux.HandleFunc("/Subcategories", ph.GetSubCategories)
	httpsServer := httptest.NewTLSServer(mux)
	defer httpsServer.Close()

	client := httpsServer.Client()
	urlString := httpsServer.URL

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie

	cookies = append(cookies, &entity.CookieMock)

	u, _ := url.Parse(urlString)
	jar.SetCookies(u, cookies)

	client.Jar = jar
	form := url.Values{}
	form.Add("category", entity.ProjectMock.Category)
	resp, err := client.PostForm(urlString+"/Subcategories", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()
}
