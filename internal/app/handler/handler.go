package handler

import (
	"lr2/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Основные маршруты
	router.GET("/", h.RedirectToMaterials)
	router.GET("/materials", h.GetAllMaterials)
	router.GET("/material/:id", h.GetMaterialByID)

	// Маршруты для заявок на датирование
	router.POST("/dating/add-material", h.AddMaterialToDatingRequest)
	router.GET("/datinganalysisrequest/:id", h.GetDatingRequest)
	router.POST("/dating/delete", h.DeleteDatingRequest)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static/styles", "./resources/styles")
	router.Static("/static/img", "./resources/img")
}

func (h *Handler) RedirectToMaterials(ctx *gin.Context) {
	ctx.Redirect(302, "/materials")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
