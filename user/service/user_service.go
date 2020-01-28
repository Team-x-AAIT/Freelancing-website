package service

import (
	"errors"
	"fmt"
	"net/smtp"
	"regexp"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/stringTools"
	"github.com/Team-x-AAIT/Freelancing-website/user"
	"golang.org/x/crypto/bcrypt"
)

// Service is a struct that defines the UserService type.
type Service struct {
	conn user.IRepository
}

// NewUserService is a function that returns a new UserService type.
func NewUserService(connection user.IRepository) *Service {
	return &Service{conn: connection}
}

// Verification is a method that verify that eligibilty of a new user.
func (service *Service) Verification(userHolder *entity.User, idnt entity.Identification) entity.ErrorBag {

	var errMap entity.ErrorBag = make(map[string]string)
	tempUser, origin := service.conn.SearchTPUser(userHolder.Email)

	if tempUser.UID != "" && origin == idnt.From {
		userHolder.Password = tempUser.Password
		return nil
	}
	if idnt.TpFlag {
		userHolder.Password = stringTools.RandomStringGN(10)
		tempUser = service.conn.SearchUser(userHolder.Email)
		if tempUser.UID != "" {
			errMap["email"] = "email already exists"
			return errMap
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userHolder.Password), 12)
		userHolder.Password = string(hashedPassword)

		service.RegisterUser(userHolder)
		service.conn.RegisterTPUsers(userHolder.UID, userHolder.Email, userHolder.Password, idnt.From)
		return nil
	}

	if idnt.From == "EditProfile" {
		errMap := service.AuthUserProfile(userHolder, idnt)
		if len(errMap) != 0 {
			return errMap
		}
		return nil
	}

	verificationToken := stringTools.RandomStringGN(30)
	tempVerificationPack := entity.TempVerificationPack{User: userHolder, Token: verificationToken}
	entity.UserTempRepository = append(entity.UserTempRepository, tempVerificationPack)
	errMap = service.AuthUserProfile(userHolder, idnt)
	if len(errMap) != 0 {
		return errMap
	}
	service.SendVerification(tempVerificationPack)
	return nil
}

// RegisterUser is a method that register a new user to the system.
func (service *Service) RegisterUser(user *entity.User) error {
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
		return errors.New("invalid")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("invalid")
	}

	return nil
}

// EditUserProfile is a method that updates the user type using the provided parameter.
func (service *Service) EditUserProfile(user *entity.User) error {

	if user.ProfilePic != "" {
		filename := service.conn.RemoveFileDB(user.UID, "profile_pic")
		if filename != "" {
			if err := service.conn.RemoveFile(filename, "profile_pic"); err != nil {
				return err
			}
		}

	} else {
		xUser := service.SearchUser(user.UID)
		user.ProfilePic = xUser.ProfilePic
	}

	if user.CV != "" {
		filename := service.conn.RemoveFileDB(user.UID, "cv")
		if filename != "" {
			if err := service.conn.RemoveFile(filename, "cv"); err != nil {
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
func (service *Service) SearchUser(identifier string) *entity.User {
	user := service.conn.SearchUser(identifier)
	if user.UID == "" {
		return nil
	}
	return user
}

// RemoveUser remove a user from user and related informations from the table.
func (service *Service) RemoveUser(uid string) (*entity.User, error) {

	userHolder, _ := service.RemoveTPUser(uid)

	matchTags := service.GetMatchTags(uid)
	for _, value := range matchTags {
		WorkType := fmt.Sprintf("%d", value.WorkType)
		err := service.RemoveMatchTag(value.UID, value.Category, value.Subcategory, WorkType)
		if err != nil {
			return nil, err
		}
	}

	userHolder, err := service.conn.RemoveUser(uid)
	if err := service.conn.RemoveFile(userHolder.ProfilePic, "profile_pic"); err != nil {
		return nil, err
	}
	if err := service.conn.RemoveFile(userHolder.CV, "cv"); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return userHolder, nil

}

// RemoveTPUser remove a user from user_tp table.
func (service *Service) RemoveTPUser(uid string) (*entity.User, error) {

	userHolder, err := service.conn.RemoveTPUser(uid)
	if err != nil {
		return nil, err
	}

	return userHolder, nil
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

// RemoveMatchTag is a method that facilitaties the removing of match tag of a user.
func (service *Service) RemoveMatchTag(uid string, category string, subcategory string, worktypes string) error {

	worktype := 0

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

	if err := service.conn.RemoveMatchTag(uid, category, subcategory, worktype); err != nil {
		return err
	}

	return nil
}

// GetMatchTags is a method that returns all the match tag of a user a slice of match tags.
func (service *Service) GetMatchTags(uid string) []*entity.MatchTag {

	matchTagStore := service.conn.GetUserMatchTags(uid)
	return matchTagStore
}

// SearchProjectWMatchTag is a method that returns all projects that match user matching tags.
func (service *Service) SearchProjectWMatchTag(uid string) []*entity.ProjectUserContainer {

	matchTagStore := service.GetMatchTags(uid)
	var filteredProjects []*entity.ProjectUserContainer

	for _, matchTag := range matchTagStore {

		projects := service.conn.SearchProjectWMatchTag(matchTag)

		for _, uniqueProject := range projects {
			unique := true
			for _, value := range filteredProjects {
				if uniqueProject.ID == value.Project.ID {
					unique = false
					break
				}
			}
			if unique {
				owner := service.conn.SearchUser(service.conn.GetOwner(uniqueProject.ID))
				if owner.UID != uid {
					projectUserContainer := new(entity.ProjectUserContainer)
					projectUserContainer.Firstname = owner.Firstname
					projectUserContainer.Lastname = owner.Lastname
					projectUserContainer.Phonenumber = owner.Phonenumber
					projectUserContainer.Email = owner.Email
					projectUserContainer.JobTitle = owner.JobTitle
					projectUserContainer.Country = owner.Country
					projectUserContainer.City = owner.City
					projectUserContainer.Gender = owner.Gender
					projectUserContainer.ProfilePic = owner.ProfilePic
					projectUserContainer.Project = uniqueProject
					filteredProjects = append(filteredProjects, projectUserContainer)
				}

			}
		}
	}

	return filteredProjects
}

// AuthUserProfile is a method that authenticate a user type has a valid content or not.
func (service *Service) AuthUserProfile(user *entity.User, idnt entity.Identification) entity.ErrorBag {

	var errMap entity.ErrorBag = make(map[string]string)
	matchFirstname, _ := regexp.MatchString(`^[a-zA-Z]\w*$`, user.Firstname)
	matchLastname, _ := regexp.MatchString(`^\w*$`, user.Lastname)
	matchEmail, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9\._\-&!?=#]*`, user.Email)
	matchPhonenumber, _ := regexp.MatchString(`^(\d{9})?$`, user.Phonenumber)
	matchPassword, _ := regexp.MatchString(`^\w{8}\w*$`, user.Password)
	matchGender, _ := regexp.MatchString(`^[MFO]?$`, user.Gender)

	if !matchFirstname {
		errMap["firstname"] = "Firstname should contain at least one character!"
	}
	if !matchLastname {
		errMap["lastname"] = "Invalid lastname!"
	}
	if !matchEmail {
		errMap["email"] = "Invalid email!"
	}
	if !matchPassword {
		errMap["password"] = "Invalid password!"
	}
	if !matchPhonenumber {
		errMap["phoneNumber"] = "Invalid phonenumber!"
	}
	if user.Password != idnt.ConfirmPassword {
		errMap["password"] = "Password doesn't match!"
	}
	if !matchGender {
		errMap["gender"] = "Invalid gender!"
	}

	if len(user.Bio) > 100 {
		errMap["bio"] = "Should be less than 100 characters!"
	}

	tempUser := service.conn.SearchUser(user.Email)

	if tempUser.UID != "" && idnt.From != "EditProfile" {
		errMap["email"] = "Email already exists!"
	}
	return errMap
}

// SendVerification is a method that take a temporary verification packet and send a verification email to provided email address.
func (service *Service) SendVerification(toBeVerified entity.TempVerificationPack) {
	auth := smtp.PlainAuth("", "biniyam.s72@gmail.com", "0911732688", "smtp.gmail.com")
	to := []string{toBeVerified.User.Email}
	msg := []byte("To:" + toBeVerified.User.Email + "\r\n" +
		"Subject: Verify your email address\r\n" + "\r\n" +
		`Verify your email address to complete registration` + "\r\n" + "\r\n" +
		`Hi ` + toBeVerified.User.Firstname + ",\r\n" + "\r\n" +
		`Thanks for your interest in Fjobs! To complete your registration, we need you to verify your email address.` + "\r\n" +
		"Open the Link to verify: " +
		`http://localhost:1234/Verify?token=` + toBeVerified.Token + "&email=" + toBeVerified.User.Email)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "biniyam.s72@gmail.com", to, msg)
	if err != nil {
		panic(err)
	}
}

// GetOwner is a method that serves the owner of a project.
func (service *Service) GetOwner(pid string) string {
	owner := service.conn.GetOwner(pid)
	return owner
}
