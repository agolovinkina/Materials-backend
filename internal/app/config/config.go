package config

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	JWT         JWTConfig
	Redis       RedisConfig
}

type JWTConfig struct {
	Token         string
	ExpiresIn     time.Duration
	SigningMethod jwt.SigningMethod
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DB          int
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	// JWT конфигурация
	cfg.JWT.Token = getEnv("JWT_SECRET", "your-secret-key")
	expiresInStr := getEnv("JWT_EXPIRES_IN", "24h")
	cfg.JWT.ExpiresIn, err = time.ParseDuration(expiresInStr)
	if err != nil {
		cfg.JWT.ExpiresIn = 24 * time.Hour
	}
	cfg.JWT.SigningMethod = jwt.SigningMethodHS256

	// Redis конфигурация
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port, _ = strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.User = getEnv("REDIS_USER", "")
	cfg.Redis.DB, _ = strconv.Atoi(getEnv("REDIS_DB", "0"))
	cfg.Redis.DialTimeout = 10 * time.Second
	cfg.Redis.ReadTimeout = 10 * time.Second

	log.Info("config parsed")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
