package entity

import (
	"time"
)

//Admin represents application admin
type Admin struct {
	AID         string
	Firstname   string
	Lastname    string
	Password    string
	Phonenumber string
	Email       string
}

// Project is a stuct that defines the type Project.
type Project struct {
	ID            string    `json: "id"`
	Title         string    `json: "title"`
	Description   string    `json: "description"`
	Details       string    `json: "details"`
	AttachedFiles []string  `json: "attached_files"`
	Category      string    `json: "category"`
	Subcategory   string    `json: "subcategory"`
	Budget        float64   `json: "budget"`
	WorkType      int64     `json: "worktype"`
	Closed        bool      `json: "closed"`
	CreatedAt     time.Time `json: "created_at"`
}

// User is a struct that defines the type User.
type User struct {
	UID         string
	Firstname   string
	Lastname    string
	Password    string
	Phonenumber string
	Email       string
	JobTitle    string
	Country     string
	City        string
	Gender      string
	CV          string
	ProfilePic  string
	Bio         string
	Prefe       int64
	Rating      float32
}

// ApplicationBag is a struct that holds the application form or result
type ApplicationBag struct {
	PID         string `json: "pid"`
	ApplicantID string `json: "uid"`
	Proposal    string `json: "proposal"`
	Hired       bool
	Seen        bool
	Status      int64
	CreatedAt   time.Time
	Project     *Project
}

// ProjectUserContainer is a struct that containe both the project and owner information.
type ProjectUserContainer struct {
	Firstname     string
	Lastname      string
	Phonenumber   string
	Email         string
	JobTitle      string
	Country       string
	City          string
	Gender        string
	ProfilePic    string
	Project       *Project
	CreatedString string
}

// MatchTag is a struct that define a type of a match tag returned from database.
type MatchTag struct {
	UID         string
	Category    string
	Subcategory string
	WorkType    int
}

// SearchBag is a struct implementation of search information that facilitates the searching process.
type SearchBag struct {
	SearchKey     string `json: "search_key"`
	SearchBy      string `json: "search_by"`
	FilterTypeS   string `json: "filter_type"`
	FilterValue1S string `json: "filter_value1"`
	FilterValue2S string `json: "filter_value2"`
	PageNumS      string `json: "page_num"`
}

// Applicants is a struct that holds the use information and his/her proposal statment.
type Applicants struct {
	ApplicantUID string    `json: "applicantid"`
	Firstname    string    `json: "firstname"`
	Lastname     string    `json: "lastname"`
	Email        string    `json: "email"`
	Gender       string    `json: "gender"`
	JobTitle     string    `json: "jobtitle"`
	PhoneNumber  string    `json: "phonenumber"`
	ProfilePic   string    `json: "profilepic"`
	Rating       float32   `json: "rating"`
	CreatedAt    time.Time `json: "createdat"`
	Proposal     string    `json: "proposal"`
	Hired        bool      `json: "hired"`
}

//Session represents login user session
type Session struct {
	ID        uint
	SID       string
	Expires   int64
	SecretKey []byte
}

// Identification is a struct that hold a basic information about third party authentication.
type Identification struct {
	TpFlag          bool
	From            string
	ConfirmPassword string
}

// TempVerificationPack is a struct that contain a user type pointer and a verification token.
type TempVerificationPack struct {
	User  *User
	Token string
}

// UserTempRepository is a slice that contain a temporary verificaion packet that are waiting to be verified.
var UserTempRepository []TempVerificationPack

// ErrorBag is a map that stores the founded validation form errors.
type ErrorBag map[string]string

// Get method to retrieve the first error message for a given// field from the map.
func (errMap ErrorBag) Get(field string) string {
	errS := errMap[field]
	return errS
}

// NewProject is a function that returns a new Project type from provided arguments.
func NewProject(title, description, details, catagory, subcatagroy string, budget float64, worktype int64) *Project {
	project := Project{
		Title:       title,
		Description: description,
		Details:     details,
		Category:    catagory,
		Subcategory: subcatagroy,
		Budget:      budget,
		WorkType:    worktype}

	return &project

}

// NewUser is a function that return a new User type from provided arguments.
func NewUser(firstname, lastname, email, profilePic, password, phonenumber, jobTitle, country, city, gender, cv, bio string) *User {
	user := User{
		Firstname:   firstname,
		Lastname:    lastname,
		ProfilePic:  profilePic,
		Password:    password,
		Phonenumber: phonenumber,
		Email:       email,
		JobTitle:    jobTitle,
		Country:     country,
		City:        city,
		Gender:      gender,
		CV:          cv,
		Bio:         bio,
	}
	return &user
}

// NewUserFR is a function that return a new User type from provided arguments.
func NewUserFR(firstname, lastname, email, password string) *User {
	user := User{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  password,
	}
	return &user
}
