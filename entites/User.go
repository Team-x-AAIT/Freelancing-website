package entites

// User is a struct that defines the type User.
type User struct {
	Firstname   string
	Lastname    string
	Phonenumber string
	Email       string
	JobTitle    string
	Password    string
	Country     string
	City        string
	Gender      string
}

// NewUser is a function that return a new User type from provided arguments.
func NewUser(firstname, lastname, phonenumber, email, jobtitle, password, country, city, gender string) *User {
	user := User{
		Firstname:   firstname,
		Lastname:    lastname,
		Phonenumber: phonenumber,
		Email:       email,
		JobTitle:    jobtitle,
		Password:    password,
		Country:     country,
		City:        city,
		Gender:      gender,
	}
	return &user
}
