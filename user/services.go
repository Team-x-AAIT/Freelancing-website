package user

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// IService is an interface that specifies what a User type can do.
type IService interface {
	RegisterUser(*entity.User) error
	Login(string, string) error
	EditUserProfile(*entity.User) error
	Verification(*entity.User, entity.Identification) entity.ErrorBag
	SearchUser(string) *entity.User
	AddMatchTag(uid string, category string, subcategory string, worktype string) error
	RemoveMatchTag(uid string, category string, subcategory string, worktype string) error
	GetMatchTags(uid string) []*entity.MatchTag
	SearchProjectWMatchTag(uid string) []*entity.ProjectUserContainer
	GetOwner(string) string

	RemoveUser(uid string) (*entity.User, error)
	RemoveTPUser(uid string) (*entity.User, error)
	// AuthUserProfile(*entity.User, Identification) error
}

// SessionService is an interface that specifie the session related operations.
type SessionService interface {
	Session(sessionID string) (*entity.Session, error)
	StoreSession(session *entity.Session) (*entity.Session, error)
	DeleteSession(sessionID string) (*entity.Session, error)
}
