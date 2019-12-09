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
	Rating      float32
}

// NewUser is a function that return a new User type from provided arguments.
func NewUser(firstname, lastname, password, phonenumber, email, jobtitle, country, city, gender string) *User {
	user := User{
		Firstname:   firstname,
		Lastname:    lastname,
		Password:    password,
		Phonenumber: phonenumber,
		Email:       email,
		JobTitle:    jobtitle,
		Country:     country,
		City:        city,
		Gender:      gender,
	}
	return &user
}
