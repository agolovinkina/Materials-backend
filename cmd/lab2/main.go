package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"lr4/internal/app/config"
	"lr4/internal/app/ds"
	"lr4/internal/app/dsn"
	"lr4/internal/app/handler"
	"lr4/internal/app/repository"
	"lr4/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
// @securityDefinitions.cookie SessionCookie
// @name session_token
// @description Используется для запросов из браузера.
// @in cookie
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

	// Инициализация Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	logrus.Info("Успешное подключение к Redis.")

	// Подключение к PostgreSQL
	postgresString := dsn.FromEnv()
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
	rep := repository.NewRepository(db)

	// Инициализация MinIO
	minioClient, err := config.NewMinioClient(conf.Minio)
	if err != nil {
		logrus.Fatal("Failed to initialize Minio client: ", err)
	}

	// Инициализация сервиса расчета календарных дат
	calcService := service.NewCalendarDateCalculator()

	// Парсинг JWT длительности
	jwtDuration, err := time.ParseDuration(conf.JWT.ExpiresIn)
	if err != nil {
		logrus.Fatalf("Invalid JWT expiration duration in config: %v", err)
	}

	// Инициализация handler
	hand := handler.NewHandler(
		rep,
		minioClient,
		conf.Minio.BucketName,
		rdb,
		conf.JWT.SecretKey,
		conf.ServiceHost,
		jwtDuration,
		calcService,
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
