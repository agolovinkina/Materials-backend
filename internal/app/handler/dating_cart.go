package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddMaterialToCartRequest struct {
	MaterialID uint `json:"material_id" binding:"required"`
}

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

func (h *Handler) GetDatingCart(ctx *gin.Context) {
	userID := h.getCurrentUserID()

	request, err := h.Repository.GetCurrentDatingRequest(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, request)
}

func (h *Handler) AddMaterialToDatingCart(ctx *gin.Context) {
	userID := h.getCurrentUserID()

	var req AddMaterialToCartRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("DEBUG: Adding material %d to cart for user %d\n", req.MaterialID, userID)

	material, err := h.Repository.GetMaterialByID(req.MaterialID)
	if err != nil {
		fmt.Printf("DEBUG: Material not found: %v\n", err)
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	fmt.Printf("DEBUG: Material found: %s\n", material.MaterialName)

	err = h.Repository.AddMaterialToDatingRequest(userID, req.MaterialID,
		"",
		"Образец для радиоуглеродного анализа",
		1.0)

	if err != nil {
		fmt.Printf("DEBUG: Error adding material to request: %v\n", err)
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

	fmt.Printf("DEBUG: Material successfully added to cart\n")

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Материал успешно добавлен в заявку",
	})
}
