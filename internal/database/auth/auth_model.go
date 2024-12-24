package auth

import "time"

type Session struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}
