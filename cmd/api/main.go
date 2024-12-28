// @title Car Social API
// @version 1.0
// @description API сервер для Car Social приложения

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	_ "github.com/NikitaBelov-mobile/car-social/docs"
	"github.com/NikitaBelov-mobile/car-social/internal/config"
	"github.com/NikitaBelov-mobile/car-social/internal/database"
	authDatabase "github.com/NikitaBelov-mobile/car-social/internal/database/auth"
	userDatabase "github.com/NikitaBelov-mobile/car-social/internal/database/user"
	"github.com/NikitaBelov-mobile/car-social/internal/service/token"
	authHandler "github.com/NikitaBelov-mobile/car-social/internal/transport/http/handler/auth"
	userHandler "github.com/NikitaBelov-mobile/car-social/internal/transport/http/handler/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	jwtService, err := token.NewTokenManager("asdasd")

	userDB := userDatabase.NewUserRepositoryImpl(db)
	authDB := authDatabase.NewAuthRepositoryImpl(db)

	userRoute := userHandler.NewHandler(userDB)
	authRoute := authHandler.NewHandler(userDB, authDB, jwtService)

	router := gin.Default()

	userRoute.Register(&router.RouterGroup)
	authRoute.Register(&router.RouterGroup)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}
