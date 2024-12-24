package user

import (
	"database/sql"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(user *User) error {
	query := `
        INSERT INTO users (phone, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $3)
        RETURNING id`

	now := time.Now()
	return r.db.QueryRow(query,
		user.Phone,
		user.PasswordHash,
		now,
	).Scan(&user.ID)
}
