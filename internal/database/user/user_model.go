package user

import "time"

type User struct {
	ID           int       `db:"id"`
	Phone        string    `db:"phone"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
