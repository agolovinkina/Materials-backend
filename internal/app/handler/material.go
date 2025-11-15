package handler

import (
	"fmt"

	"net/http"
	"strconv"

	"lr4/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CreateMaterialRequest struct {
	MaterialName        string `json:"material_name" binding:"required"`
	MaterialDescription string `json:"material_description"`
	Isotopes            string `json:"isotopes" binding:"required"`
	SampleRequirements  string `json:"sample_requirements" binding:"required"`
}

type UpdateMaterialRequest struct {
	MaterialName        string `json:"material_name"`
	MaterialDescription string `json:"material_description"`
	Isotopes            string `json:"isotopes"`
	SampleRequirements  string `json:"sample_requirements"`
}

// GetAllMaterials получает все материалы
// @Summary Get all materials
// @Description Get all materials with optional search
// @Tags Materials
// @Produce json
// @Param search query string false "Search query"
// @Success 200 {object} map[string]interface{} "Success response with materials array"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /materials [get]
func (h *Handler) GetAllMaterials(ctx *gin.Context) {
	searchQuery := ctx.Query("search")

	var materials []ds.Material
	var err error

	if searchQuery == "" {
		materials, err = h.Repository.GetAllMaterials()
	} else {
		materials, err = h.Repository.SearchMaterialsByName(searchQuery)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, materials)
}

// GetMaterialByID получает материал по ID
// @Summary Get material by ID
// @Description Get material details by ID
// @Tags Materials
// @Produce json
// @Param id path int true "Material ID"
// @Success 200 {object} map[string]interface{} "Success response with material"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Material not found"
// @Router /materials/{id} [get]
func (h *Handler) GetMaterialByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	h.successResponse(ctx, material)
}

// CreateMaterial создает новый материал
// @Summary Create material
// @Description Create new material (admin only)
// @Tags Materials
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param material body CreateMaterialRequest true "Material data"
// @Success 201 {object} map[string]interface{} "Material created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /materials [post]
func (h *Handler) CreateMaterial(ctx *gin.Context) {
	var req CreateMaterialRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	material := &ds.Material{
		MaterialName:        req.MaterialName,
		MaterialDescription: req.MaterialDescription,
		Isotopes:            req.Isotopes,
		SampleRequirements:  req.SampleRequirements,
		IsDeleted:           false,
	}

	err := h.Repository.CreateMaterial(material)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, material)
}

// UpdateMaterial обновляет материал
// @Summary Update material
// @Description Update material by ID (admin only)
// @Tags Materials
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Material ID"
// @Param material body UpdateMaterialRequest true "Material update data"
// @Success 200 {object} map[string]interface{} "Material updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /materials/{id} [put]
func (h *Handler) UpdateMaterial(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req UpdateMaterialRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем, существует ли материал
	_, err = h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("материал не найден"))
		return
	}

	updates := make(map[string]interface{})
	if req.MaterialName != "" {
		updates["material_name"] = req.MaterialName
	}
	if req.MaterialDescription != "" {
		updates["material_description"] = req.MaterialDescription
	}
	if req.Isotopes != "" {
		updates["isotopes"] = req.Isotopes
	}
	if req.SampleRequirements != "" {
		updates["sample_requirements"] = req.SampleRequirements
	}

	if len(updates) == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нет данных для обновления"))
		return
	}

	err = h.Repository.UpdateMaterial(uint(id), updates)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получаем обновленный материал для ответа
	updatedMaterial, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, updatedMaterial)
}

// DeleteMaterial удаляет материал (мягкое удаление)
// @Summary Delete material
// @Description Delete material by ID (admin only)
// @Tags Materials
// @Security BearerAuth
// @Param id path int true "Material ID"
// @Success 200 {object} map[string]interface{} "Material deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /materials/{id} [delete]
func (h *Handler) DeleteMaterial(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем, существует ли материал
	_, err = h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("материал не найден"))
		return
	}

	err = h.Repository.DeleteMaterial(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Материал успешно удален",
	})
}

// UploadMaterialImage
// @Summary Загрузить изображение материала
// @Description Загружает и обновляет изображение для материала по ID. Требуются права **Модератора**.
// @Tags Материалы
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID материала"
// @Param image formData file true "Файл изображения"
// @Success 200 {object} ds.Material "Успешная загрузка, возвращает обновленный материал"
// @Failure 400 {object} handler.ErrorResponse "Ошибка загрузки/формата файла"
// @Failure 403 {object} handler.ErrorResponse "Доступ запрещен (не модератор)"
// @Failure 500 {object} handler.ErrorResponse "Ошибка Minio/сервера"
// @Router /materials/{id}/image [post]
func (h *Handler) UploadMaterialImage(ctx *gin.Context) {
	if h.MinioClient == nil {
		h.errorHandler(ctx, http.StatusServiceUnavailable,
			fmt.Errorf("image storage service not configured"))
		return
	}

	if !h.checkMinioConnection() {
		h.errorHandler(ctx, http.StatusServiceUnavailable,
			fmt.Errorf("image storage service temporarily unavailable. Please try again later"))
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid material ID: %v", err))
		return
	}

	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("material not found: %v", err))
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("image file is required: %v", err))
		return
	}

	if !h.isValidImage(file) {
		h.errorHandler(ctx, http.StatusBadRequest,
			fmt.Errorf("invalid image format. Allowed: JPEG, PNG, GIF, WebP"))
		return
	}

	if file.Size > 5*1024*1024 {
		h.errorHandler(ctx, http.StatusBadRequest,
			fmt.Errorf("image size too large. Maximum 5MB allowed"))
		return
	}

	if material.MaterialImageURL != "" {
		err = h.deleteImageFromMinio(material.MaterialImageURL)
		if err != nil {
			logrus.Warnf("Failed to delete old image from Minio: %v", err)
		}
	}

	objectName := h.generateImageName(file.Filename, id)

	imageURL, err := h.uploadImageToMinio(file, objectName)
	if err != nil {
		logrus.Errorf("Failed to upload image to Minio: %v", err)
		h.errorHandler(ctx, http.StatusInternalServerError,
			fmt.Errorf("failed to upload image: %v", err))
		return
	}

	updatedMaterial, err := h.Repository.UpdateMaterialImage(uint(id), imageURL)
	if err != nil {
		logrus.Errorf("Failed to update material in database, rolling back image upload")
		h.deleteImageFromMinio(imageURL)
		h.errorHandler(ctx, http.StatusInternalServerError,
			fmt.Errorf("failed to update material: %v", err))
		return
	}

	logrus.Infof("Successfully updated material %d with new image: %s", id, imageURL)

	ctx.JSON(http.StatusOK, gin.H{
		"data": updatedMaterial,
	})
}
