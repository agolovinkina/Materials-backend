package config

import "os"

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	SSL       bool
}

func NewMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		Bucket:    getEnv("MINIO_BUCKET", "materials"),
		SSL:       getEnv("MINIO_SSL", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
