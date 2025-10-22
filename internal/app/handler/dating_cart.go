package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddMaterialToCartRequest struct {
	MaterialID uint `json:"material_id" binding:"required"`
}

// GetDatingCartIcon - метод 7
func (h *Handler) GetDatingCartIcon(ctx *gin.Context) {
	userID := h.getCurrentUserID()

	count, err := h.Repository.GetDatingCartCount(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	response := gin.H{
		"count": count,
	}

	h.successResponse(ctx, response)
}

// GetDatingCart - новый метод для получения черновика заявки
func (h *Handler) GetDatingCart(ctx *gin.Context) {
	userID := h.getCurrentUserID()

	request, err := h.Repository.GetCurrentDatingRequest(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, request)
}

// AddMaterialToDatingCart - метод 21
func (h *Handler) AddMaterialToDatingCart(ctx *gin.Context) {
	userID := h.getCurrentUserID()

	var req AddMaterialToCartRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем материал для расчета вероятности
	material, err := h.Repository.GetMaterialByID(req.MaterialID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Рассчитываем вероятность на основе характеристик материала
	probability := h.CalcService.CalculateProbability(1.0, material.Isotopes, material.SampleRequirements)

	err = h.Repository.AddMaterialToDatingRequest(userID, req.MaterialID,
		"", // пустой комментарий
		probability,
		"Образец для радиоуглеродного анализа",
		1.0) // стандартный вес

	if err != nil {
		if err.Error() == "материал уже добавлен в заявку" {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Материал уже добавлен в заявку",
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
