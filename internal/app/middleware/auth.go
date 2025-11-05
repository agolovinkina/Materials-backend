package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"lr4/internal/app/config"
	"lr4/internal/app/ds"
	"lr4/internal/app/redis"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthMiddleware struct {
	cfg         *config.Config
	redisClient *redis.Client
}

func NewAuthMiddleware(cfg *config.Config, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

// GetJWTSecret возвращает секретный ключ для JWT
func (a *AuthMiddleware) GetJWTSecret() string {
	return a.cfg.JWT.Token
}

// GetJWTExpiresIn возвращает время жизни JWT токена
func (a *AuthMiddleware) GetJWTExpiresIn() time.Duration {
	return a.cfg.JWT.ExpiresIn
}

func (a *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Authorization header required",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Invalid authorization format",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Проверяем, не в blacklist ли токен (если Redis доступен)
		if a.redisClient != nil {
			err := a.redisClient.CheckJWTInBlacklist(ctx.Request.Context(), tokenString)
			if err == nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error":   true,
					"message": "Token revoked",
				})
				ctx.Abort()
				return
			}
			if !errors.Is(err, redis.Nil) {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": "Internal server error",
				})
				ctx.Abort()
				return
			}
		}

		token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.cfg.JWT.Token), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Invalid token",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*ds.JWTClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		// Сохраняем user_id в контекст
		ctx.Set("user_id", claims.Subject)
		ctx.Set("user_uuid", claims.UserUUID.String())
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}
