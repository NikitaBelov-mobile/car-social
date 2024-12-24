package user

import (
	"database/sql"
	"errors"
	"time"
)

type UserRepository interface {
	Create(user *User) error
	GetByPhone(phone string) (*User, error)
	GetByID(id int) (*User, error)
	Update(user *User) error
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *User) error {
	query := `
        INSERT INTO users (phone, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $3)
        RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query,
		user.Phone,
		user.PasswordHash,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) GetByPhone(phone string) (*User, error) {
	user := &User{}
	query := `
        SELECT id, phone, password_hash, created_at, updated_at
        FROM users
        WHERE phone = $1`

	err := r.db.QueryRow(query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByID(id int) (*User, error) {
	user := &User{}
	query := `
        SELECT id, phone, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(user *User) error {
	query := `
        UPDATE users
        SET phone = $1,
            password_hash = $2,
            updated_at = $3
        WHERE id = $4
        RETURNING created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query,
		user.Phone,
		user.PasswordHash,
		now,
		user.ID,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}
