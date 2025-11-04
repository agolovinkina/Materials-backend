package main

import (
	"fmt"

	"lr2/internal/app/config"
	"lr2/internal/app/dsn"
	"lr2/internal/app/handler"
	"lr2/internal/app/repository"
	"lr2/internal/app/service"
	"lr2/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println("PostgreSQL DSN:", postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	// Инициализация MinIO
	minioConfig := config.NewMinIOConfig()
	minioService, err := service.NewMinIOService(
		minioConfig.Endpoint,
		minioConfig.AccessKey,
		minioConfig.SecretKey,
		minioConfig.Bucket,
		minioConfig.SSL,
	)
	if err != nil {
		logrus.Fatalf("error initializing MinIO: %v", err)
	}
	logrus.Info("MinIO service initialized successfully")

	// Инициализация сервиса расчета календарных дат
	calcService := service.NewCalendarDateCalculator()

	hand := handler.NewHandler(rep, calcService, minioService)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
