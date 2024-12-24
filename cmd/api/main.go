package main

import (
	"log"

	"github.com/NikitaBelov-mobile/car-social/internal/config"
	"github.com/NikitaBelov-mobile/car-social/internal/database"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Получаем экземпляр sql.DB для проверки соединения
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Проверяем соединение
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
}
