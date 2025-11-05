package middleware

import (
	"lr4/internal/app/role"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminRequired проверяет что пользователь администратор
func (a *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Доступ запрещен",
			})
			ctx.Abort()
			return
		}

		// Получаем роль как значение типа role.Role
		userRoleValue, ok := userRole.(role.Role)
		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Доступ запрещен: неверный формат роли",
			})
			ctx.Abort()
			return
		}

		// Сравниваем с константой role.Admin
		if userRoleValue != role.Admin {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Требуются права администратора",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// ModeratorRequired проверяет что пользователь модератор или админ
func (a *AuthMiddleware) ModeratorRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Доступ запрещен",
			})
			ctx.Abort()
			return
		}

		// Получаем роль как значение типа role.Role
		userRoleValue, ok := userRole.(role.Role)
		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Доступ запрещен: неверный формат роли",
			})
			ctx.Abort()
			return
		}

		// Сравниваем с константами role.Manager и role.Admin
		if userRoleValue != role.Manager && userRoleValue != role.Admin {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Требуются права модератора",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
