package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
<a href="/LinkedInLogin">Log in with linked in</a>
<a href="/FacebookLogin">Log in with Facebook</a>
</body></html>
`

var (

	// URepositoryDB is a pointer to the user PsqlUserRepository type.
	URepositoryDB IRepository
	// UService is a pointer to the user Service type.
	UService IService

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
	// UserTempRepository is a slice that contain a temporary verificaion packet that are waiting to be verified.
	UserTempRepository []TempVerificationPack
)

// HandleMain is a Handler func that parse the main window.
func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

// HandleGoogleLogin is a Handler func that initiate google oauth process to the google.Endpoint.
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	OauthStateString = randomStringGN(20)
	url := googleOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback is a Handler func that handle all the token verification and profile information from google api.
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
	ServerAUT = randomStringGN(20)
	postQuery := "thirdParty=true" +
		"&firstname=" + tempUser.GivenName +
		"&lastname=" + tempUser.FamilyName +
		"&email=" + tempUser.Email +
		"&from=google" +
		"&serverAUT=" + ServerAUT
	_, err = http.Post("http://localhost:1234/Register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(postQuery)))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// HandleFacebookLogin is a Handler func that initiate google oauth process to the facebook.Endpoint.
func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	OauthStateString = randomStringGN(20)
	url := facebookOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleFacebookCallback is a Handler func that handle all the token verification and profile information from facebook api.
func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {

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

	ServerAUT = randomStringGN(20)
	postQuery := "thirdParty=true" +
		"&firstname=" + tempUser.FirstName +
		"&lastname=" + tempUser.LastName +
		"&email=" + tempUser.Email +
		"&from=facebook" +
		"&serverAUT=" + ServerAUT

	_, err = http.Post("http://localhost:1234/Register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(postQuery)))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// HandleLinkedInLogin is a Handler func that initiate google oauth process to the linkedin.Endpoint.
func HandleLinkedInLogin(w http.ResponseWriter, r *http.Request) {
	OauthStateString = randomStringGN(20)
	url := linkedinOauthConfig.AuthCodeURL(OauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleLinkedInCallback is a Handler func that handle all the token verification and profile information from linkedin api.
func HandleLinkedInCallback(w http.ResponseWriter, r *http.Request) {

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

	ServerAUT = randomStringGN(20)
	postQuery := "thirdParty=true" +
		"&firstname=" + tempUser.FirstName +
		"&lastname=" + tempUser.LastName +
		"&email=" + tempUser.Email +
		"&from=linkedin" +
		"&serverAUT=" + ServerAUT
	_, err = http.Post("http://localhost:1234/Register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(postQuery)))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// Register is a Handler func that initaite registration process.
func Register(w http.ResponseWriter, r *http.Request) {

	var userHolder *entities.User
	var password string

	thirdParty := r.FormValue("thirdParty")
	var identification Identification
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

	userHolder = entities.NewUserFR(firstname, lastname, email, password)

	err := UService.Verification(userHolder, identification)
	if err != nil {
		panic(err)
	}
	if identification.TpFlag {
		// until session is being used
		http.Redirect(w, r, "/Login?email="+userHolder.Email+"&password="+userHolder.Password, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/Check_Your_Email", http.StatusSeeOther)
	}
	return

}

// Verify is a Handler func that verify a token from a request query and user to database if valid.
func Verify(w http.ResponseWriter, r *http.Request) {
	verificationToken := r.URL.Query().Get("token")
	email := r.URL.Query().Get("email")

	for _, value := range UserTempRepository {
		if value.token == verificationToken && value.user.Email == email {
			err := UService.RegisterUser(value.user)
			if err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/Login?email="+value.user.Email+"&password="+value.user.Password, http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "invalid verification Token!", http.StatusBadRequest)
}

// Index is check function.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sucessfull!"))
}
