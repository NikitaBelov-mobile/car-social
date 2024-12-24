package auth

import (
	"database/sql"
	"errors"
	"time"
)

type AuthRepository interface {
	CreateSession(session *Session) error
	GetSessionByRefreshToken(refreshToken string) (*Session, error)
	DeleteSession(refreshToken string) error
	DeleteUserSessions(userID int) error
}

type AuthRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepositoryImpl(db *sql.DB) AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) CreateSession(session *Session) error {
	query := `
        INSERT INTO sessions (user_id, refresh_token, expires_at, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	return r.db.QueryRow(
		query,
		session.UserID,
		session.RefreshToken,
		session.ExpiresAt,
		time.Now(),
	).Scan(&session.ID)
}

func (r *AuthRepositoryImpl) GetSessionByRefreshToken(refreshToken string) (*Session, error) {
	session := &Session{}
	query := `
        SELECT id, user_id, refresh_token, expires_at, created_at
        FROM sessions
        WHERE refresh_token = $1`

	err := r.db.QueryRow(query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	return session, nil
}

func (r *AuthRepositoryImpl) DeleteSession(refreshToken string) error {
	query := `DELETE FROM sessions WHERE refresh_token = $1`
	result, err := r.db.Exec(query, refreshToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}

func (r *AuthRepositoryImpl) DeleteUserSessions(userID int) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
