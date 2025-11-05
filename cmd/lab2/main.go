package main

import (
	"context"
	"fmt"
	"log"

	"lr4/internal/app/config"
	"lr4/internal/app/ds"
	"lr4/internal/app/dsn"
	"lr4/internal/app/handler"
	"lr4/internal/app/middleware"
	"lr4/internal/app/redis"
	"lr4/internal/app/repository"
	"lr4/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "lr4/cmd/docs"
)

// @title Система анализа материалов
// @version 1.0
// @description Бэкенд сервис для анализа материалов и радиоуглеродного датирования.
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Используется для запросов через Insomnia/Postman: "Bearer <JWT>"
func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// Подключение к PostgreSQL
	postgresString := dsn.FromEnv()
	if postgresString == "" {
		log.Fatal("DSN string is empty. Check your environment variables")
	}

	fmt.Println("PostgreSQL DSN:", postgresString)

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(postgresString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database!")

	// Выполняем миграции
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Material{},
		&ds.MaterialAnalysisRequest{},
		&ds.RequestMaterial{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully!")

	// Инициализация репозитория
	rep := &repository.Repository{DB: db}

	// Инициализация Redis клиента
	var redisClient *redis.Client
	redisClient, err = redis.New(context.Background(), conf.Redis)
	if err != nil {
		logrus.Warnf("Redis not available: %v", err)
	} else {
		logrus.Info("Успешное подключение к Redis.")
	}

	// Инициализация сервиса расчета календарных дат
	calcService := service.NewCalendarDateCalculator()

	// Инициализация middleware аутентификации
	authMiddleware := middleware.NewAuthMiddleware(conf, redisClient)

	// Инициализация MinIO (опционально)
	var minioService *service.MinIOService
	minioConfig := config.NewMinIOConfig()
	minioService, err = service.NewMinIOService(
		minioConfig.Endpoint,
		minioConfig.AccessKey,
		minioConfig.SecretKey,
		minioConfig.Bucket,
		minioConfig.SSL,
	)
	if err != nil {
		logrus.Warnf("MinIO not available: %v", err)
	} else {
		logrus.Info("MinIO service initialized successfully")
	}

	hand := handler.NewHandler(
		rep,
		calcService,
		minioService,
		authMiddleware,
		redisClient,
	)

	// Регистрация маршрутов
	hand.RegisterHandler(router)
	hand.RegisterStatic(router)

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
