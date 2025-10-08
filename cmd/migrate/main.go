package main

import (
	"log"

	"lr2/internal/app/ds"
	"lr2/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Получаем DSN строку
	dsnString := dsn.FromEnv()
	if dsnString == "" {
		log.Fatal("DSN string is empty. Check your environment variables")
	}

	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
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
}
