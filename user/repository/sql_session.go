package repository

import (
	"database/sql"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
	"github.com/Team-x-AAIT/Freelancing-website/user"
)

// SessionRepository is a struct that defines the SessionRepository type.
type SessionRepository struct {
	conn *sql.DB
}

// NewSessionRepository returns a new SessionRepository type.
func NewSessionRepository(db *sql.DB) user.SessionRepository {
	return &SessionRepository{conn: db}
}

// Session returns a given stored session
func (psql *SessionRepository) Session(sessionID string) (*entity.Session, error) {
	session := entity.Session{}
	secretkey := "b"

	stmt, err := psql.conn.Prepare("SELECT * FROM Sessions WHERE session_id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(sessionID)

	row.Scan(
		&session.ID,
		&session.SID,
		&session.Expires,
		&secretkey)
	session.SecretKey = []byte(secretkey)

	return &session, err
}

// StoreSession stores a given session
func (psql *SessionRepository) StoreSession(session *entity.Session) (*entity.Session, error) {

	sess := session
	stmt, _ := psql.conn.Prepare(`INSERT INTO Sessions (id, session_id, expires, secretkey) VALUES (?,?,?,?)`)
	_, err := stmt.Exec(
		session.ID,
		session.SID,
		session.Expires,
		string(session.SecretKey))

	if err != nil {
		return nil, err
	}
	return sess, err
}

// DeleteSession deletes a given session
func (psql *SessionRepository) DeleteSession(sessionID string) (*entity.Session, error) {
	sess, err := psql.Session(sessionID)
	if err != nil {
		return nil, err
	}

	stmt, _ := psql.conn.Prepare(`DELETE FROM Sessions WHERE session_id=?`)
	_, err = stmt.Exec(sessionID)

	if err != nil {
		return nil, err
	}
	return sess, err
}
