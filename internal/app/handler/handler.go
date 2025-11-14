package handler

import (
	"lr4/internal/app/middleware"
	"lr4/internal/app/repository"
	"lr4/internal/app/service"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
	BucketName  string
	RedisClient *redis.Client
	SecretKey   string
	HostName    string
	JWTDur      time.Duration
	CalcService *service.CalendarDateCalculator
}

func NewHandler(r *repository.Repository, minioClient *minio.Client, bucketName string, rdb *redis.Client, secretKey, hostName string, jwtDur time.Duration, calcService *service.CalendarDateCalculator) *Handler {
	return &Handler{
		Repository:  r,
		MinioClient: minioClient,
		BucketName:  bucketName,
		RedisClient: rdb,
		SecretKey:   secretKey,
		HostName:    hostName,
		JWTDur:      jwtDur,
		CalcService: calcService,
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static/styles", "./resources/styles")
	router.Static("/static/img", "./resources/img")
	router.Static("/static/uploads", "./uploads")
}

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

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// getCurrentUserIDFromContext получает ID пользователя из контекста
func (h *Handler) getCurrentUserIDFromContext(ctx *gin.Context) int {
	return middleware.GetUserID(ctx)
}

// IsAdmin проверяет, является ли пользователь администратором/модератором
func (h *Handler) IsAdmin(ctx *gin.Context) bool {
	userRole := middleware.GetUserRole(ctx)
	return userRole == "moderator"
}

// getCurrentModeratorID получает ID модератора из контекста
func (h *Handler) getCurrentModeratorID() int {
	return middleware.GetUserID(nil) // Будет получать из контекста в реальном использовании
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/auth/register", h.RegisterUser)
		api.POST("/auth/login", h.LoginUser)

		api.GET("/materials", h.GetAllMaterials)
		api.GET("/materials/:id", h.GetMaterialByID)

		api.GET("/dating/cart/icon", h.GetDatingCartIcon)

		api.POST("/test-upload", h.TestFileUpload)
	}

	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware(h.SecretKey, h.RedisClient))
	{
		auth.POST("/auth/logout", h.LogoutUser)
		auth.GET("/user/profile", h.GetProfile)
		auth.PUT("/user/profile", h.UpdateProfile)

		auth.GET("/dating/cart", h.GetDatingCart)
		auth.POST("/dating/cart/materials", h.AddMaterialToDatingCart)

		auth.GET("/dating/requests", h.GetDatingRequests)
		auth.GET("/dating/requests/:id", h.GetDatingRequestByID)
		auth.PUT("/dating/requests/:id", h.UpdateDatingRequest)
		auth.PUT("/dating/requests/:id/form", h.FormDatingRequest)
		auth.DELETE("/dating/requests/:id", h.DeleteDatingRequest)

		auth.DELETE("/dating/request-materials/:requestId/:materialId", h.RemoveMaterialFromRequest)
		auth.PUT("/dating/request-materials/:requestId/:materialId", h.UpdateRequestMaterial)
	}

	moderator := auth.Group("")
	moderator.Use(middleware.RequireModerator())
	{
		moderator.POST("/materials", h.CreateMaterial)
		moderator.PUT("/materials/:id", h.UpdateMaterial)
		moderator.DELETE("/materials/:id", h.DeleteMaterial)
		moderator.POST("/materials/:id/image", h.UploadMaterialImage)

		moderator.PUT("/dating/requests/:id/process", h.ProcessDatingRequest)
	}
}
