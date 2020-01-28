package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/session"
	"github.com/Team-x-AAIT/Freelancing-website/stringTools"
	"github.com/Team-x-AAIT/Freelancing-website/user"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

// UserHandler is a struct for the user hanlding functions.
type UserHandler struct {
	UService user.IService
	SService user.SessionService
	Temp     *template.Template
	CSRF     []byte
}

// InputContainer is a struct that holds information for the templates.
type InputContainer struct {
	LoggedInUser *entity.User
	Error        entity.ErrorBag
	Counter      [3]int64
	Prefe        bool
	CSRF         string
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:1234/GoogleCallback",
		ClientID:     "866699224887-chqgletm3pgv85d8t1j3hqu9qbggg9c6.apps.googleusercontent.com",
		ClientSecret: "EcGwbY1ikvFKs2kPhKSQPk8X",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}

	linkedinOauthConfig = &oauth2.Config{
		ClientID:     "86uo5nbpfdo1ac",
		ClientSecret: "bJSiYAeW4XhhlaTB",
		RedirectURL:  "http://localhost:1234/LinkedInCallback",
		Scopes:       []string{"r_liteprofile", "r_emailaddress"},
		Endpoint:     linkedin.Endpoint,
	}

	facebookOauthConfig = &oauth2.Config{
		ClientID:     "453733495330808",
		ClientSecret: "4b5856af660b4440654e8c4eb31c3a45",
		RedirectURL:  "http://localhost:1234/FacebookCallback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	// ServerAUT is a sceret key that is used to identify the POST request, if it checks out then the call was made from the server it self.
	ServerAUT string
	// OauthStateString is a sceret state key that is used for identifying a request in api exchange.
	OauthStateString string
)

// NewUserHandler is a function that return new User Handler type.
func NewUserHandler(service user.IService, sservice user.SessionService, temp *template.Template, csrf []byte) *UserHandler {
	return &UserHandler{UService: service, SService: sservice, Temp: temp, CSRF: csrf}
}

// WelcomePage is a Hanler func that handles the main
func (uh *UserHandler) WelcomePage(w http.ResponseWriter, r *http.Request) {
	uh.Temp.ExecuteTemplate(w, "WelcomePage.html", nil)
}

// HandleGoogleLogin is a Handler func that initiate google oauth process to the google.Endpoint.
func (uh *UserHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	OauthStateString = stringTools.RandomStringGN(20)
	url := googleOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// HandleGoogleCallback is a Handler func that handle all the token verification and profile information from google api.
func (uh *UserHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	statCheck := r.FormValue("state")
	if statCheck != OauthStateString {
		http.Error(w, fmt.Sprintf("Wrong state string: Expected %s, got %s. Please, try again", OauthStateString, statCheck), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Code exchange failed with %s", err.Error()), http.StatusBadRequest)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	var tempUser struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Profile       string `json:"profile"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Gender        string `json:"gender"`
	}
	err = json.Unmarshal(contents, &tempUser)
	if err != nil {
		panic(err)
	}
	ServerAUT = stringTools.RandomStringGN(20)
	postQuery := "/Register?thirdParty=true" +
		"&firstname=" + tempUser.GivenName +
		"&lastname=" + tempUser.FamilyName +
		"&email=" + tempUser.Email +
		"&from=google" +
		"&serverAUT=" + ServerAUT
	http.Redirect(w, r, postQuery, http.StatusSeeOther)

}

// HandleFacebookLogin is a Handler func that initiate google oauth process to the facebook.Endpoint.
func (uh *UserHandler) HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	OauthStateString = stringTools.RandomStringGN(20)
	url := facebookOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// HandleFacebookCallback is a Handler func that handle all the token verification and profile information from facebook api.
func (uh *UserHandler) HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {

	statCheck := r.FormValue("state")
	if statCheck != OauthStateString {
		http.Error(w, fmt.Sprintf("Wrong state string: Expected %s, got %s. Please, try again", OauthStateString, statCheck), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := facebookOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Code exchange failed with %s", err.Error()), http.StatusBadRequest)
		return
	}

	response, err := http.Get(`https://graph.facebook.com/me?access_token=` + token.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	var checkToken struct {
		ID   string `json: "id"`
		Name string `json:"name"`
	}
	err = json.Unmarshal(contents, &checkToken)
	if err != nil {
		panic(err)
	}
	response, err = http.Get(`https://graph.facebook.com/` + checkToken.ID + `?fields=id,name,email,first_name,last_name&access_token=` + token.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	contents, err = ioutil.ReadAll(response.Body)

	var tempUser struct {
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	err = json.Unmarshal(contents, &tempUser)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(tempUser.Email))

	ServerAUT = stringTools.RandomStringGN(20)
	postQuery := "/Register?thirdParty=true" +
		"&firstname=" + tempUser.FirstName +
		"&lastname=" + tempUser.LastName +
		"&email=" + tempUser.Email +
		"&from=facebook" +
		"&serverAUT=" + ServerAUT

	http.Redirect(w, r, postQuery, http.StatusSeeOther)

}

// HandleLinkedInLogin is a Handler func that initiate google oauth process to the linkedin.Endpoint.
func (uh *UserHandler) HandleLinkedInLogin(w http.ResponseWriter, r *http.Request) {
	OauthStateString = stringTools.RandomStringGN(20)
	url := linkedinOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// HandleLinkedInCallback is a Handler func that handle all the token verification and profile information from linkedin api.
func (uh *UserHandler) HandleLinkedInCallback(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	statCheck := r.FormValue("state")
	if OauthStateString != statCheck {
		http.Error(w, fmt.Sprintf("Wrong state string: Expected %s, got %s. Please, try again", OauthStateString, statCheck), http.StatusBadRequest)
		return
	}

	token, err := linkedinOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := linkedinOauthConfig.Client(oauth2.NoContext, token)
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Set("Bearer", token.AccessToken)
	response, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer response.Body.Close()
	str, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tempUser struct {
		ID        string
		FirstName string `json:"localizedFirstName"`
		LastName  string `json:"localizedLastName"`
		Email     string `json:"emailAddress"`
		Pp        struct {
			Picture string `json:"displayImage"`
		} `json:"profilePicture"`
	}

	err = json.Unmarshal(str, &tempUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ServerAUT = stringTools.RandomStringGN(20)
	postQuery := "/Register?thirdParty=true" +
		"&firstname=" + tempUser.FirstName +
		"&lastname=" + tempUser.LastName +
		"&email=" + tempUser.Email +
		"&from=linkedin" +
		"&serverAUT=" + ServerAUT
	http.Redirect(w, r, postQuery, http.StatusSeeOther)

}

// Register is a Handler func that initaite registration process.
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	var userHolder *entity.User
	var password string

	if r.Method == http.MethodGet {
		uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(uh.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		inputContainer := InputContainer{CSRF: token}
		uh.Temp.ExecuteTemplate(w, "SignUp.html", inputContainer)
		return
	}

	if r.Method == http.MethodPost {

		thirdParty := r.FormValue("thirdParty")
		var identification entity.Identification
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		identification.ConfirmPassword = r.FormValue("confirmPassword")

		if thirdParty == "true" {

			if r.FormValue("serverAUT") != ServerAUT {
				http.Error(w, "Invalid server key", http.StatusBadRequest)
				return
			}

			identification.From = r.FormValue("from")
			identification.TpFlag = true
		} else {
			password = r.FormValue("password")
			identification.ConfirmPassword = r.FormValue("confirmPassword")
		}

		// Validating CSRF Token
		csrfToken := r.FormValue("csrf")
		ok, errCRFS := stringTools.ValidCSRF(csrfToken, uh.CSRF)

		userHolder = entity.NewUserFR(firstname, lastname, email, password)
		errMap := uh.UService.Verification(userHolder, identification)
		if !ok || errCRFS != nil {
			if len(errMap) == 0 {
				errMap = make(map[string]string)
			}
			errMap["csrf"] = "Invalid token used!"
		}
		if len(errMap) > 0 {
			uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, _ := stringTools.CSRFToken(uh.CSRF)
			inputContainer := InputContainer{Error: errMap, CSRF: token}
			uh.Temp.ExecuteTemplate(w, "SignUp.html", inputContainer)
			return
		}

		if identification.TpFlag {

			newSession := uh.configSess()
			claims := stringTools.Claims(email, newSession.Expires)
			session.Create(claims, newSession, w)
			_, err := uh.SService.StoreSession(newSession)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/Dashboard", http.StatusSeeOther)
		}

		uh.Temp.ExecuteTemplate(w, "CheckEmail.html", nil)
		return
	}
}

// Login is a Handler func that initaite login process.
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	errMap := make(map[string]string)
	user := uh.Authentication(r)
	if user != nil {
		http.Redirect(w, r, "/Dashboard", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(uh.CSRF)
		inputContainer := InputContainer{CSRF: token}
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		uh.Temp.ExecuteTemplate(w, "LoginPage.html", inputContainer)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validating CSRF Token
		csrfToken := r.FormValue("csrf")
		ok, errCRFS := stringTools.ValidCSRF(csrfToken, uh.CSRF)
		if !ok || errCRFS != nil {
			errMap["csrf"] = "Invalid token used!"
		}

		err := uh.UService.Login(email, password)

		if err != nil || len(errMap) > 0 {
			errMap["login"] = "Invalid email or password!"
			uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, _ := stringTools.CSRFToken(uh.CSRF)
			inputContainer := InputContainer{Error: errMap, CSRF: token}
			uh.Temp.ExecuteTemplate(w, "LoginPage.html", inputContainer)
			return
		}

		newSession := uh.configSess()
		claims := stringTools.Claims(email, newSession.Expires)
		session.Create(claims, newSession, w)
		_, err = uh.SService.StoreSession(newSession)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/Dashboard", http.StatusSeeOther)
		return
	}

}

// UpdateProfile is a Handler func that initaite the updating profile process.
func (uh *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	loggedInUser := uh.Authentication(r)
	if loggedInUser == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
		token, err := stringTools.CSRFToken(uh.CSRF)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		inputContainer := InputContainer{LoggedInUser: loggedInUser, CSRF: token}
		uh.Temp.ExecuteTemplate(w, "ProfilePage.html", inputContainer)
		return
	}

	if r.Method == http.MethodPost {

		var profilePic, cv string
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		phonenumber := r.FormValue("phonenumber")
		email := r.FormValue("email")
		jobTitle := r.FormValue("jobTitle")
		country := r.FormValue("country")
		city := r.FormValue("city")
		gender := r.FormValue("gender")
		bio := r.FormValue("bio")
		prefeS := r.FormValue("prefe")
		prefe, _ := strconv.ParseInt(prefeS, 0, 0)

		// In Edit profile password will not be change be still it will be authenticated!
		password := "ValidPassword123"

		csrfToken := r.FormValue("csrf")
		ok, errCRFS := stringTools.ValidCSRF(csrfToken, uh.CSRF)

		user := entity.NewUser(firstname, lastname, email, profilePic, password, phonenumber, jobTitle, country, city, gender, cv, bio)
		user.UID = loggedInUser.UID
		user.Prefe = prefe

		errMap := uh.UService.Verification(user, entity.Identification{ConfirmPassword: "ValidPassword123", From: "EditProfile"})
		if !ok || errCRFS != nil {
			if len(errMap) == 0 {
				errMap = make(map[string]string)
			}
			errMap["csrf"] = "Invalid token used!"
		}

		if len(errMap) > 0 {
			uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, _ := stringTools.CSRFToken(uh.CSRF)
			inputContainer := InputContainer{LoggedInUser: user, Error: errMap, CSRF: token}
			uh.Temp.ExecuteTemplate(w, "ProfilePage.html", inputContainer)
			return
		}

		pFilename, err := uh.ResourceExtractorUser(user.UID, "profilePic", "image", r)

		if err != nil && (err.Error() == "file to large" || err.Error() == "invalid format") {
			if len(errMap) == 0 {
				errMap = make(map[string]string)
			}
			errMap["profilePic"] = "File size should be less than 5MB!"
			uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, _ := stringTools.CSRFToken(uh.CSRF)
			inputContainer := InputContainer{LoggedInUser: user, Error: errMap, CSRF: token}
			uh.Temp.ExecuteTemplate(w, "ProfilePage.html", inputContainer)
			return
		}

		cFilename, err := uh.ResourceExtractorUser(user.UID, "cv", "file", r)

		if err != nil && (err.Error() == "file too large" || err.Error() == "invalid format") {
			if len(errMap) == 0 {
				errMap = make(map[string]string)
			}
			switch err.Error() {
			case "invalid format":
				errMap["cv"] = "Only pdf format is allowed!"
			case "file too large":
				errMap["cv"] = "File size should be less than 5MB!"
			}

			uh.CSRF, _ = stringTools.GenerateRandomBytes(30)
			token, _ := stringTools.CSRFToken(uh.CSRF)
			inputContainer := InputContainer{LoggedInUser: user, Error: errMap, CSRF: token}
			uh.Temp.ExecuteTemplate(w, "ProfilePage.html", inputContainer)
			return
		}

		user.ProfilePic = pFilename
		user.CV = cFilename
		err = uh.UService.EditUserProfile(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
}

// Verify is a Handler func that verify a token from a request query and user to database if valid.
func (uh *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	verificationToken := r.URL.Query().Get("token")
	email := r.URL.Query().Get("email")

	for _, value := range entity.UserTempRepository {
		if value.Token == verificationToken && value.User.Email == email {

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(value.User.Password), 12)
			value.User.Password = string(hashedPassword)

			err := uh.UService.RegisterUser(value.User)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			newSession := uh.configSess()
			claims := stringTools.Claims(email, newSession.Expires)
			session.Create(claims, newSession, w)
			_, err = uh.SService.StoreSession(newSession)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/Dashboard", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "invalid verification Token!", http.StatusBadRequest)
}

// Authentication is a function that checks if the user has already logged in or not.
func (uh *UserHandler) Authentication(r *http.Request) *entity.User {

	newSession := uh.configSess()
	fjCookie, err := r.Cookie(newSession.SID)
	if err != nil {
		return nil
	}
	value := fjCookie.Value

	newSession, err = uh.SService.Session(value)

	if err != nil {
		return nil
	}

	if newSession.Expires < time.Now().Unix() {
		uh.SService.DeleteSession(value)
		return nil
	}

	ok, _ := session.Valid(value, newSession.SecretKey)
	if !ok {
		return nil
	}

	email, _ := stringTools.GenerateValue(value, newSession.SecretKey)

	user := uh.UService.SearchUser(email)
	return user
}

func (uh *UserHandler) configSess() *entity.Session {
	tokenExpires := time.Now().Add(time.Hour * 240).Unix()
	sessionID := "Fjobs_User_Cookie"
	signingString := stringTools.RandomStringGN(32)

	signingKey := []byte(signingString)

	return &entity.Session{
		Expires:   tokenExpires,
		SecretKey: signingKey,
		SID:       sessionID,
	}
}

// Logout is a Handler func that perform logout operation by revoking authentication pass.
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	newSession := uh.configSess()
	cookie, _ := r.Cookie(newSession.SID)

	session.Remove(newSession.SID, w)
	uh.SService.DeleteSession(cookie.Value)
	http.Redirect(w, r, "/Login", http.StatusSeeOther)
}

// AddMatchTag is a Handler func that accepts the matching tags store them.
func (uh *UserHandler) AddMatchTag(w http.ResponseWriter, r *http.Request) {

	user := uh.Authentication(r)
	if user == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	uid := user.UID
	category := r.FormValue("category")
	subcategory := r.FormValue("subcategory")
	worktype := r.FormValue("worktype")

	if err := uh.UService.AddMatchTag(uid, category, subcategory, worktype); err != nil {
		panic(err)
	}
	w.Write([]byte("okay"))

}

// RemoveMatchTag is a Handler func that accepts the matching tags and remove it.
func (uh *UserHandler) RemoveMatchTag(w http.ResponseWriter, r *http.Request) {

	user := uh.Authentication(r)
	if user == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	uid := user.UID
	category := r.FormValue("category")
	subcategory := r.FormValue("subcategory")
	worktype := r.FormValue("worktype")

	if err := uh.UService.RemoveMatchTag(uid, category, subcategory, worktype); err != nil {
		panic(err)
	}

	w.Write([]byte("okay"))
}

// GetMatchTags is a Handler func that sends all the match tag the user have.
func (uh *UserHandler) GetMatchTags(w http.ResponseWriter, r *http.Request) {
	user := uh.Authentication(r)
	if user == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	matchTagStore := uh.UService.GetMatchTags(user.UID)
	json, err := json.Marshal(matchTagStore)
	if err != nil {
		panic(err)
	}

	w.Write(json)

}

// GetProjectsWMatchTags is a Handler func that sends all the projects taht match the user match Tags.
func (uh *UserHandler) GetProjectsWMatchTags(w http.ResponseWriter, r *http.Request) {

	user := uh.Authentication(r)
	if user == nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	uid := user.UID

	projects := uh.UService.SearchProjectWMatchTag(uid)

	json, err := json.Marshal(projects)
	if err != nil {
		panic(err)
	}

	w.Write(json)
}

// ResourceExtractorUser is a method that extract file from a request.
func (uh *UserHandler) ResourceExtractorUser(uid string, name string, fileType string, r *http.Request) (string, error) {

	count := stringTools.RandomStringGN(10)

	fm, fh, err := r.FormFile(name)

	if err != nil {
		return "", err
	}

	tempFile, _ := ioutil.ReadAll(fm)
	tempFileType := http.DetectContentType(tempFile)

	if tempFileType != "application/pdf" && !strings.HasPrefix(tempFileType, "image") {
		return "", errors.New("invalid format")
	}

	if fh.Size > 5000000 {
		return "", errors.New("file too large")
	}
	defer fm.Close()

	path, _ := os.Getwd()
	suffix := ""
	endPoint := 0

	for i := len(fh.Filename) - 1; i >= 0; i-- {
		if fh.Filename[i] == '.' {
			endPoint = i
			break
		}
	}

	for ; endPoint < len(fh.Filename); endPoint++ {
		suffix += string(fh.Filename[endPoint])
	}

	NewFileName := fmt.Sprintf("asset_"+fileType+"_"+uid+"_%s"+suffix, count)
	if name == "profilePic" {
		column := "profile_pic"
		path = filepath.Join(path, "..", "..", "ui", "assets", column, NewFileName)
	} else {
		path = filepath.Join(path, "..", "..", "ui", "assets", name, NewFileName)
	}

	out, _ := os.Create(path)
	defer out.Close()

	_, err = io.Copy(out, fm)

	if err != nil {
		panic(err)
	}

	return NewFileName, nil

}

// Index is check function.
func (uh *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	uh.Temp.ExecuteTemplate(w, "CheckEmail.html", nil)
}

// IndexNot is check function.
func (uh *UserHandler) IndexNot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Sucessfull!"))
}
