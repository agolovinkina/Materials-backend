package handler

import (
	"lr2/internal/app/repository"
	"lr2/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	CalcService *service.CalendarDateCalculator
}

func NewHandler(r *repository.Repository, calcService *service.CalendarDateCalculator) *Handler {
	return &Handler{
		Repository:  r,
		CalcService: calcService,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Домен услуги (6 методов)
		api.GET("/materials", h.GetAllMaterials)
		api.GET("/materials/:id", h.GetMaterialByID)
		api.POST("/materials", h.CreateMaterial)
		api.PUT("/materials/:id", h.UpdateMaterial)
		api.DELETE("/materials/:id", h.DeleteMaterial)
		api.POST("/materials/:id/image", h.UploadMaterialImage)

		// Домен заявки
		api.GET("/dating/cart/icon", h.GetDatingCartIcon)
		api.GET("/dating/cart", h.GetDatingCart)
		api.GET("/dating/requests", h.GetDatingRequests)

		// Группа для операций с заявками
		datingGroup := api.Group("/dating")
		{
			// Основные операции с заявками
			datingGroup.GET("/requests/:id", h.GetDatingRequestByID)
			datingGroup.PUT("/requests/:id", h.UpdateDatingRequest)
			datingGroup.PUT("/requests/:id/form", h.FormDatingRequest)
			datingGroup.PUT("/requests/:id/process", h.ProcessDatingRequest)
			datingGroup.DELETE("/requests/:id", h.DeleteDatingRequest)

			// Операции с материалами в заявках
			datingGroup.DELETE("/request-materials/:requestId/:materialId", h.RemoveMaterialFromRequest)
			datingGroup.PUT("/request-materials/:requestId/:materialId", h.UpdateRequestMaterial)

			// Добавление в заявку
			datingGroup.POST("/cart/materials", h.AddMaterialToDatingCart)
		}

		// Домен пользователь (5 методов)
		api.POST("/user/register", h.RegisterUser)
		api.GET("/user/profile", h.GetProfile)
		api.PUT("/user/profile", h.UpdateProfile)
		api.POST("/user/login", h.LoginUser)
		api.POST("/user/logout", h.LogoutUser)
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static/styles", "./resources/styles")
	router.Static("/static/img", "./resources/img")
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

// Временная функция для имитации авторизации
func (h *Handler) getCurrentUserID() int {
	return 1
}

func (h *Handler) getCurrentModeratorID() uint {
	return 2
}
