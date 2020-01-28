package repository

import (
	"database/sql"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

// MockSessionRepository is a struct that defines the SessionRepository type.
type MockSessionRepository struct {
	conn *sql.DB
}

// NewMockSessionRepository returns a new MockSessionRepository type.
func NewMockSessionRepository(db *sql.DB) *MockSessionRepository {
	return &MockSessionRepository{conn: db}
}

// Session returns a given stored session
func (psql *MockSessionRepository) Session(sessionID string) (*entity.Session, error) {
	return &entity.SessionMock, nil
}

// StoreSession stores a given session
func (psql *MockSessionRepository) StoreSession(session *entity.Session) (*entity.Session, error) {
	session.ID = entity.SessionMock.ID
	return &entity.SessionMock, nil
}

// DeleteSession deletes a given session
func (psql *MockSessionRepository) DeleteSession(sessionID string) (*entity.Session, error) {
	return &entity.SessionMock, nil
}
