package user

import (
	"errors"
	"math/rand"
	"net/smtp"
	"regexp"
	"time"

	"github.com/Team-x-AAIT/Freelancing-website/entities"
)

// Service is a struct that defines the UserService type.
type Service struct {
	conn IRepository
}

// Identification is a struct that hold a basic information about third party authentication.
type Identification struct {
	TpFlag          bool
	From            string
	ConfirmPassword string
}

// TempVerificationPack is a struct that contain a user type pointer and a verification token.
type TempVerificationPack struct {
	user  *entities.User
	token string
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
	Project       *entities.Project
	CreatedString string
}

// IService is an interface that specifies what a User type can do.
type IService interface {
	RegisterUser(*entities.User) error
	Login(string, string) error
	EditUserProfile(*entities.User) error
	Verification(*entities.User, Identification) error
	SearchUser(string) *entities.User
	AddMatchTag(uid string, category string, subcategory string, worktype string) error
	RemoveMatchTag(uid string, category string, subcategory string, worktype string) error
	GetMatchTags(uid string) []*MatchTag
	SearchProjectWMatchTag(uid string) []*ProjectUserContainer
	// AuthUserProfile(*entities.User, Identification) error
}

// NewService is a function that returns a new UserService type.
func NewService(connection IRepository) IService {
	return &Service{conn: connection}
}

// Verification is a method that verify that eligibilty of a new user.
func (service *Service) Verification(userHolder *entities.User, idnt Identification) error {

	tempUser, origin := service.conn.SearchTPUser(userHolder.Email)

	if tempUser.UID != "" && origin == idnt.From {
		userHolder.Password = tempUser.Password
		return nil
	}
	if idnt.TpFlag {
		userHolder.Password = randomStringGN(10)
		tempUser = service.conn.SearchUser(userHolder.Email)
		if tempUser.UID != "" {
			return errors.New("email already exists")
		}
		service.RegisterUser(userHolder)
		service.conn.RegisterTPUsers(userHolder.UID, userHolder.Email, userHolder.Password, idnt.From)
		return nil
	}

	if idnt.From == "EditProfile" {
		err := service.AuthUserProfile(userHolder, idnt)
		if err != nil {
			return err
		}
		return nil
	}

	verificationToken := randomStringGN(30)
	tempVerificationPack := TempVerificationPack{user: userHolder, token: verificationToken}
	UserTempRepository = append(UserTempRepository, tempVerificationPack)
	err := service.AuthUserProfile(userHolder, idnt)
	if err != nil {
		return err
	}
	service.SendVerification(tempVerificationPack)
	return nil
}

// RegisterUser is a method that register a new user to the system.
func (service *Service) RegisterUser(user *entities.User) error {
	err := service.conn.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

// Login is a method that validate a user using Email and Password.
func (service *Service) Login(email string, password string) error {

	user := service.conn.SearchUser(email)
	if user.UID == "" {
		return errors.New("please use a valid email addres")
	}

	if user.Password != password {
		return errors.New("incorrect password")
	}

	return nil
}

// EditUserProfile is a method that updates the user type using the provided parameter.
func (service *Service) EditUserProfile(user *entities.User) error {

	if user.ProfilePic != "" {
		filename := URepositoryDB.RemoveFileDB(user.UID, "profile_pic")
		if filename != "" {
			if err := URepositoryDB.RemoveFile(filename, "profile_pic"); err != nil {
				return err
			}
		}

	} else {
		xUser := service.SearchUser(user.UID)
		user.ProfilePic = xUser.ProfilePic
	}

	if user.CV != "" {
		filename := URepositoryDB.RemoveFileDB(user.UID, "cv")
		if filename != "" {
			if err := URepositoryDB.RemoveFile(filename, "cv"); err != nil {
				return err
			}
		}
	} else {
		xUser := service.SearchUser(user.UID)
		user.CV = xUser.CV
	}

	err := service.conn.UpdateUser(user)

	if err != nil {
		return err
	}

	tempUser, _ := service.conn.SearchTPUser(user.UID)
	if tempUser.UID != "" {
		err := service.conn.UpdateTPUsers(user.UID, user.Email)
		if err != nil {
			return err
		}
	}

	return nil
}

// SearchUser is a method that is used for searching a user from the system.
func (service *Service) SearchUser(identifier string) *entities.User {
	user := service.conn.SearchUser(identifier)
	if user.UID == "" {
		return nil
	}
	return user
}

// AddMatchTag is a method that facilitaties the adding process of match tag.
func (service *Service) AddMatchTag(uid string, category string, subcategory string, worktypes string) error {

	worktype := 0

	if !service.conn.SearchMember("categories", category) {
		return errors.New("unknown category")
	}

	if subcategory != "" && !service.conn.SearchMember("subcategories", subcategory) {
		return errors.New("unknown subcategory")
	}

	switch worktypes {
	case "Fixed":
		worktype = 1
	case "Perhour":
		worktype = 2
	case "Negotiable":
		worktype = 3
	case "":
		worktype = 4
	default:
		return errors.New("unknow work type")
	}

	matchTagStore := service.conn.GetUserMatchTags(uid)
	for _, value := range matchTagStore {
		if value.Category == category && value.Subcategory == subcategory && value.WorkType == worktype {
			return errors.New("duplicate tag")
		}
	}
	if len(matchTagStore) > 2 {
		return errors.New("reached maximum tag support")
	}

	err := service.conn.AddMatchTag(uid, category, subcategory, worktype)
	if err != nil {
		return err
	}

	return nil
}

