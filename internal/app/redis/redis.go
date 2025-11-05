package redis

import (
	"context"
	"errors"
	"fmt"
	"lr4/internal/app/config"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "lr4."

type Client struct {
	cfg    config.RedisConfig
	client *redis.Client
}

var Nil = errors.New("redis: nil")

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	client := &Client{}
	client.cfg = cfg

	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.Password,
		Username:    cfg.User,
		Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	client.client = redisClient

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// JWT blacklist methods
const jwtPrefix = "jwt."

func getJWTKey(token string) string {
	return servicePrefix + jwtPrefix + token
}

func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, getJWTKey(jwtStr), true, jwtTTL).Err()
}

func (c *Client) CheckJWTInBlacklist(ctx context.Context, jwtStr string) error {
	err := c.client.Get(ctx, getJWTKey(jwtStr)).Err()
	if errors.Is(err, redis.Nil) {
		return Nil // Токен не в блеклисте
	}
	return err // Ошибка или токен в блеклисте
}
