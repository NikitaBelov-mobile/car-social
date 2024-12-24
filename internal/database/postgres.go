package database

import (
	"database/sql"
	"fmt"

	"github.com/NikitaBelov-mobile/car-social/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.DBName,
		cfg.DB.Port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return db, nil
}
