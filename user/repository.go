package user

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// IRepository is an interface that specifies database operations on User type.
type IRepository interface {
	AddUser(*entity.User) error
	SearchUser(string) *entity.User
	UpdateUser(*entity.User) error
	RegisterTPUsers(string, string, string, string) error
	UpdateTPUsers(string, string) error
	SearchTPUser(string) (*entity.User, string)

	RemoveUser(uid string) (*entity.User, error)
	RemoveTPUser(uid string) (*entity.User, error)

	CountMember(string) int
	RemoveFileDB(string, string) string
	RemoveFile(string, string) error
	AddMatchTag(uid string, category string, subcategory string, worktype int) error
	RemoveMatchTag(uid string, category string, subcategory string, worktype int) error
	GetUserMatchTags(uid string) []*entity.MatchTag
	SearchProjectWMatchTag(matchTag *entity.MatchTag) []*entity.Project
	SearchMember(tableName string, columnValue string) bool
	GetOwner(string) string
}

// SessionRepository is an interface that specifie session related database operations.
type SessionRepository interface {
	Session(sessionID string) (*entity.Session, error)
	StoreSession(session *entity.Session) (*entity.Session, error)
	DeleteSession(sessionID string) (*entity.Session, error)
}
