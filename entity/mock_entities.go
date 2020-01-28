package entity

import (
	"net/http"
	"time"
)

// UserMock mocks User entity.
var UserMock = User{
	UID:         "UID1",
	Firstname:   "Benyam",
	Lastname:    "Simayehu",
	Password:    "0911732688Biny",
	Phonenumber: "0900010197",
	Email:       "binysimayehu@gmail.com",
	JobTitle:    "Software Enginner",
	Country:     "Ethiopia",
	City:        "Addis Ababa",
	Gender:      "M",
	CV:          "",
	ProfilePic:  "",
	Bio:         "",
	Prefe:       0,
	Rating:      3.5,
}

// AdminMock mocks Admin entity.
var AdminMock = Admin{
	AID:         "AID1",
	Firstname:   "Liben",
	Lastname:    "Hailu",
	Password:    "123456",
	Phonenumber: "0900010197",
	Email:       "Libenhailu@gmail.com",
}

// ProjectMock mocks the Project entity.
var ProjectMock = Project{
	ID:          "PID1",
	Title:       "I need a professional react developer.",
	Description: "I am highly in need of a professional react developer with 3 year experience.",
	Details: `Lorem ipsum dolor sit amet consectetur adipisicing elit. Commodi, nihil quia tempora
					 quam aliquam impedit molestiae enim provident odit labore minima vel itaque id autem 
					 ullam exercitationem, culpa temporibus blanditiis.`,
	AttachedFiles: []string{"assets_project_1.pdf"},
	Category:      "Web Developement",
	Subcategory:   "Front-end",
	Budget:        30,
	WorkType:      2,
	Closed:        false,
	CreatedAt:     time.Now(),
}

// ApplicationMock mocks the application entity.
var ApplicationMock = ApplicationBag{
	PID:         "PID1",
	ApplicantID: "UID1",
	Proposal: `Lorem ipsum dolor sit amet consectetur adipisicing elit. Commodi, nihil quia tempora
					 quam aliquam impedit molestiae enim provident odit labore minima vel itaque id autem 
					 ullam exercitationem, culpa temporibus blanditiis.`,
	Hired:     false,
	Seen:      false,
	Status:    1,
	CreatedAt: time.Now(),
	Project:   &ProjectMock,
}

// PUCMock mocks the ProjectUserContainer struct.
var PUCMock = ProjectUserContainer{
	Firstname:     "Benyam",
	Lastname:      "Simayehu",
	Phonenumber:   "0900010197",
	Email:         "binysimayehu@gmail.com",
	JobTitle:      "Software Enginner",
	Country:       "Ethiopia",
	Gender:        "M",
	Project:       &ProjectMock,
	CreatedString: time.Now().String(),
}

// MatchTagMock mocks the MatchTag struct.
var MatchTagMock = MatchTag{
	UID:         "UID1",
	Category:    "Mobile Development",
	Subcategory: "Front-End",
	WorkType:    1,
}

// ApplicantsMock mocks the Applicants struct.
var ApplicantsMock = Applicants{
	ApplicantUID: "UID1",
	Firstname:    "Benyam",
	Lastname:     "Simayehu",
	Email:        "binysimayehu@gmail.com",
	JobTitle:     "Software Enginner",
	PhoneNumber:  "911732688",
	Gender:       "M",
	Rating:       3.5,
	Proposal: `Lorem ipsum dolor sit amet consectetur adipisicing elit. Commodi, nihil quia tempora
					 quam aliquam impedit molestiae enim provident odit labore minima vel itaque id autem 
					 ullam exercitationem, culpa temporibus blanditiis.`,
	Hired: false,
}

// SessionMock mocks sessions of loged in user
var SessionMock = Session{
	ID:        1,
	SID:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJpbnlzaW1heWVodUBnbWFpbC5jb20ifQ.aOYW75cKYZTWmVLG24PW81FWnA8dvR6PUGOzXZZejj0",
	SecretKey: []byte("Secret_key"),
	Expires:   time.Now().Add(time.Hour * 400).Unix(),
}

// CookieMock mocks a cookie struct.
var CookieMock = http.Cookie{
	Name:     "Fjobs_User_Cookie",
	Value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJpbnlzaW1heWVodUBnbWFpbC5jb20ifQ.aOYW75cKYZTWmVLG24PW81FWnA8dvR6PUGOzXZZejj0",
	Expires:  time.Now().Add(time.Hour * 400),
	HttpOnly: true,
}

// claims := CustomClaims{
// 		Email: "binysimayehu@gmail.com",
// 	}

// IdentificationMock mocks the Identification struct.
var IdentificationMock = Identification{
	TpFlag:          false,
	From:            "",
	ConfirmPassword: "091173bkNZ",
}
