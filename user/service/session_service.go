package service

import (
	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/user"
)

// SessionService is a struct that defines the SessionService type.
type SessionService struct {
	conn user.SessionRepository
}

// NewSessionService returns a new SessionService type.
func NewSessionService(connection user.SessionRepository) user.SessionService {
	return &SessionService{conn: connection}
}

// Session returns a given stored session
func (service *SessionService) Session(sessionID string) (*entity.Session, error) {
	sess, err := service.conn.Session(sessionID)
	if err != nil {
		return nil, err
	}
	return sess, err
}

// StoreSession stores a given session
func (service *SessionService) StoreSession(session *entity.Session) (*entity.Session, error) {
	sess, err := service.conn.StoreSession(session)
	if err != nil {
		return nil, err
	}
	return sess, err
}

// DeleteSession deletes a given session
func (service *SessionService) DeleteSession(sessionID string) (*entity.Session, error) {
	sess, err := service.conn.DeleteSession(sessionID)
	if err != nil {
		return nil, err
	}
	return sess, err
}
