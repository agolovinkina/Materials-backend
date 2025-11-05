package handler

import (
	"fmt"
	"lr4/internal/app/middleware"
	"lr4/internal/app/redis"
	"lr4/internal/app/repository"
	"lr4/internal/app/role"
	"lr4/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler структура обработчиков
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

type Handler struct {
	Repository     *repository.Repository
	CalcService    *service.CalendarDateCalculator
	MinIOService   *service.MinIOService
	AuthMiddleware *middleware.AuthMiddleware
	RedisClient    *redis.Client
}

func NewHandler(r *repository.Repository, calcService *service.CalendarDateCalculator, minioService *service.MinIOService, authMiddleware *middleware.AuthMiddleware, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:     r,
		CalcService:    calcService,
		MinIOService:   minioService,
		AuthMiddleware: authMiddleware,
		RedisClient:    redisClient,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Публичные эндпоинты
		api.POST("/auth/register", h.RegisterUser)
		api.POST("/auth/login", h.LoginUser)

		// Материалы (публичные)
		api.GET("/materials", h.GetAllMaterials)
		api.GET("/materials/:id", h.GetMaterialByID)

		// Защищенные эндпоинты
		protected := api.Group("")
		protected.Use(h.AuthMiddleware.AuthRequired())
		{
			// Пользователь
			protected.GET("/user/profile", h.GetProfile)
			protected.PUT("/user/profile", h.UpdateProfile)
			protected.POST("/user/logout", h.LogoutUser)

			// Корзина и заявки (доступны всем авторизованным пользователям)
			protected.GET("/dating/cart/icon", h.GetDatingCartIcon)
			protected.GET("/dating/cart", h.GetDatingCart)
			protected.POST("/dating/cart/materials", h.AddMaterialToDatingCart)
			protected.GET("/dating/requests", h.GetDatingRequests) // ← ПЕРЕМЕЩЕНO СЮДА
			protected.GET("/dating/requests/:id", h.GetDatingRequestByID)
			protected.PUT("/dating/requests/:id/form", h.FormDatingRequest)
			protected.PUT("/dating/requests/:id", h.UpdateDatingRequest)
			protected.DELETE("/dating/request-materials/:requestId/:materialId", h.RemoveMaterialFromRequest)
			protected.PUT("/dating/request-materials/:requestId/:materialId", h.UpdateRequestMaterial)

			// Административные эндпоинты - ТОЛЬКО ДЛЯ АДМИНОВ
			admin := protected.Group("")
			admin.Use(h.AuthMiddleware.AdminRequired()) // ✅ Проверка прав администратора
			{
				admin.POST("/materials", h.CreateMaterial)
				admin.PUT("/materials/:id", h.UpdateMaterial)
				admin.DELETE("/materials/:id", h.DeleteMaterial)
				admin.POST("/materials/:id/image", h.UploadMaterialImage)

			}

			// Эндпоинты модератора - ДЛЯ МОДЕРАТОРОВ И АДМИНОВ
			moderator := protected.Group("")
			moderator.Use(h.AuthMiddleware.ModeratorRequired()) // ✅ Проверка прав модератора
			{
				moderator.PUT("/dating/requests/:id/process", h.ProcessDatingRequest)
				moderator.DELETE("/dating/requests/:id", h.DeleteDatingRequest)
			}
		}

		// Тестовый endpoint
		api.POST("/test-upload", h.TestFileUpload)
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static/styles", "./resources/styles")
	router.Static("/static/img", "./resources/img")
	router.Static("/static/uploads", "./uploads")
}

// ==============================
// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ ДЛЯ РОЛЕЙ
// ==============================

// GetCurrentUserRole возвращает роль текущего пользователя
func (h *Handler) GetCurrentUserRole(ctx *gin.Context) role.Role {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		return role.Buyer
	}

	roleValue, ok := userRole.(role.Role)
	if !ok {
		return role.Buyer
	}

	return roleValue
}

// IsAdmin проверяет является ли пользователь администратором
func (h *Handler) IsAdmin(ctx *gin.Context) bool {
	return h.GetCurrentUserRole(ctx) == role.Admin
}

// IsModerator проверяет является ли пользователь модератором или администратором
func (h *Handler) IsModerator(ctx *gin.Context) bool {
	userRole := h.GetCurrentUserRole(ctx)
	return userRole == role.Manager || userRole == role.Admin
}

// ==============================
// ОСНОВНЫЕ МЕТОДЫ ОБРАБОТКИ
// ==============================

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"message": err.Error(),
	})
}

func (h *Handler) successResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"data": data,
	})
}

func (h *Handler) errorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"error":   true,
		"message": message,
	})
}

// Получение ID пользователя из контекста JWT
func (h *Handler) getCurrentUserIDFromContext(ctx *gin.Context) int {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}

	if id, ok := userID.(string); ok {
		var result int
		_, err := fmt.Sscanf(id, "%d", &result)
		if err != nil {
			return 0
		}
		return result
	}
	return 0
}

// Совместимость со старым кодом
func (h *Handler) getCurrentUserID() int {
	return 1
}

func (h *Handler) getCurrentModeratorID() uint {
	return 2
}
