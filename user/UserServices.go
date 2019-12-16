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

// IService is an interface that specifies what a User type can do.
type IService interface {
	RegisterUser(*entities.User) error
	Verification(*entities.User, Identification) error
	SearchUser(string) *entities.User
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

// SearchUser is a method that is used for searching a user from the system.
func (service *Service) SearchUser(identifier string) *entities.User {
	user := service.conn.SearchUser(identifier)
	if user.UID == "" {
		return nil
	}
	return user
}

// AuthUserProfile is a method that authenticate a user type has a valid content or not.
func (service *Service) AuthUserProfile(user *entities.User, idnt Identification) error {

	matchFirstname, _ := regexp.MatchString(`^[a-zA-Z]\w*$`, user.Firstname)
	matchLastname, _ := regexp.MatchString(`^\w*$`, user.Lastname)
	matchEmail, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9\._\-&!?=#]*`, user.Email)
	matchPhonenumber, _ := regexp.MatchString(`^(\+\d{12})?$`, user.Phonenumber)
	matchPassword, _ := regexp.MatchString(`^\w{8}\w*$`, user.Password)
	matchGender, _ := regexp.MatchString(`^[MFO]?$`, user.Gender)

	switch {
	case !matchFirstname:
		return errors.New("invalid firstname")
	case !matchLastname:
		return errors.New("invalid lastname")
	case !matchEmail:
		return errors.New("invalid email")
	case !matchPassword:
		return errors.New("invalid password")
	case !matchPhonenumber:
		return errors.New("invalid phoneNumber")
	case user.Password != idnt.ConfirmPassword:
		return errors.New("password doesn't match")
	case !matchGender:
		return errors.New("invalid gender")
	}

	tempUser := service.conn.SearchUser(user.Email)

	if tempUser.UID != "" && idnt.From != "EditProfile" {
		return errors.New("email already exists")
	}

	return nil
}

// SendVerification is a method that take a temporary verification packet and send a verification email to provided email address.
func (service *Service) SendVerification(toBeVerified TempVerificationPack) {
	auth := smtp.PlainAuth("", "biniyam.s72@gmail.com", "0911732688", "smtp.gmail.com")
	to := []string{toBeVerified.user.Email}
	msg := []byte("To:" + toBeVerified.user.Email + "\r\n" +
		"Subject: Verify your email address\r\n" + "\r\n" +
		`Verify your email address to complete registration` + "\r\n" + "\r\n" +
		`Hi ` + toBeVerified.user.Firstname + ",\r\n" + "\r\n" +
		`Thanks for your interest in Fjobs! To complete your registration, we need you to verify your email address.` + "\r\n" +
		"Open the Link to verify: " +
		`http://localhost:1234/Verify?token=` + toBeVerified.token + "&email=" + toBeVerified.user.Email)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "biniyam.s72@gmail.com", to, msg)
	if err != nil {
		panic(err)
	}
}

func randomStringGN(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
