package main

import (
	"fmt"
	"log"

	"github.com/NikitaBelov-mobile/car-social/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Устанавливаем режим Gin
	gin.SetMode(cfg.Server.Mode)

	// Инициализируем роутер
	router := gin.Default()

	// Простой тестовый эндпоинт
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Запускаем сервер
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	log.Fatal(router.Run(serverAddr))
}
