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
		Endpoint:  getEnvMinIO("MINIO_ENDPOINT", "localhost:9000"),
		AccessKey: getEnvMinIO("MINIO_ACCESS_KEY", "minioadmin"),
		SecretKey: getEnvMinIO("MINIO_SECRET_KEY", "minioadmin"),
		Bucket:    getEnvMinIO("MINIO_BUCKET", "materials"),
		SSL:       getEnvMinIO("MINIO_SSL", "false") == "true",
	}
}

// Переименуем функцию чтобы избежать конфликта
func getEnvMinIO(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
