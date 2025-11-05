package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddMaterialToCartRequest struct {
	MaterialID uint `json:"material_id" binding:"required"`
}

// GetDatingCartIcon получает количество материалов в корзине
// @Summary Get cart item count
// @Description Get the number of materials in dating cart
// @Tags Dating Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Success response with count"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/cart/icon [get]
func (h *Handler) GetDatingCartIcon(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("пользователь не авторизован"))
		return
	}

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

// GetDatingCart получает текущую корзину заявок
// @Summary Get dating cart
// @Description Get current user's dating cart with materials
// @Tags Dating Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Success response with cart data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/cart [get]
func (h *Handler) GetDatingCart(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("пользователь не авторизован"))
		return
	}

	request, err := h.Repository.GetCurrentDatingRequest(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, request)
}

// AddMaterialToDatingCart добавляет материал в корзину
// @Summary Add material to cart
// @Description Add material to dating cart
// @Tags Dating Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param material body AddMaterialToCartRequest true "Material data"
// @Success 201 {object} map[string]interface{} "Material added successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Material not found"
// @Failure 409 {object} map[string]interface{} "Material already in cart"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/cart/materials [post]
func (h *Handler) AddMaterialToDatingCart(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("пользователь не авторизован"))
		return
	}

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
