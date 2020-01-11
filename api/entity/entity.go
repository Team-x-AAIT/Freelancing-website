package entity

import (
	"time"
)

// user struct
type User struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	FirstName string `json:"firstname" gorm:"varchar(255); not null"`
	LastName  string `json:"lastname" gorm:"varchar(255) not null"`
	UserName  string `json:"username" gorm:"varchar(255);not null"`
	Email     string `json:"email" gorm:"varchar(255);not null;unique"`
	Password  string `json:"password" gorm:"varchar(255);not null`
	AboutYou  string `json:"aboutyou" gorm:"text"`
	Jobs      []MyJob
	Country   string `json:"country" gorm:"varchar(255); not null"`
	CreatedAt time.Time
}

// myjob struct
type MyJob struct {
	Job    string `json:"myjob"`
	User   User   `gorm:"foreignkey:UserID"`
	UserID uint   `json:"userid"`
}

// job struct
type Job struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title" gorm:"varchar(255); not null"`
	Description string `json:"description" gorm:"varchar(255); not null"`
	Category    string `json:"category" gorm:"varchar(255);  not null"`
	User        User   `gorm:"foreignkey:UserID"`
	UserID      uint   `json:"userid"`
	CreatedAt   time.Time
}

// apply struct
type Apply struct {
	User     User   `gorm:"foreignkey:UserID"`
	UserID   uint   `json:"userid"`
	Job      Job    `gorm:"foreignkey:JobID"`
	JobID    uint   `json:"jobid"`
	Proposal string `json:"proposal" gorm:"varchar(255); notnull"`
}
