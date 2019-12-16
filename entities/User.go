package entities

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
	Rating      float32
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
