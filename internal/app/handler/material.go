package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"lr2/internal/app/ds"

	"github.com/gin-gonic/gin"
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

// GetAllMaterials - метод 1
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

// GetMaterialByID - метод 2
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

// CreateMaterial - метод 3
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

// UpdateMaterial - метод 4
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

	err = h.Repository.UpdateMaterial(uint(id), updates)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteMaterial - метод 5
func (h *Handler) DeleteMaterial(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteMaterial(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UploadMaterialImage - метод 6
func (h *Handler) UploadMaterialImage(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем файл из формы
	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("ошибка получения файла: %w", err))
		return
	}

	// Проверяем тип файла
	if !isImageFile(file) {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("файл должен быть изображением (JPEG, PNG, GIF)"))
		return
	}

	// Генерируем уникальное имя файла
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("material_%d_%d%s", id, time.Now().Unix(), fileExt)

	// Сохраняем файл (временная реализация - сохраняем локально)
	uploadPath := "./uploads/materials/"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	filePath := filepath.Join(uploadPath, fileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Формируем URL (в продакшене нужно использовать CDN/Minio)
	imageURL := fmt.Sprintf("/static/uploads/materials/%s", fileName)

	// Обновляем материал в БД
	err = h.Repository.UpdateMaterialImage(uint(id), imageURL)
	if err != nil {
		// Удаляем загруженный файл в случае ошибки
		os.Remove(filePath)
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"image_url": imageURL,
		"message":   "Изображение успешно загружено",
	})
}

// Вспомогательная функция для проверки типа файла
func isImageFile(file *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	fileHeader, err := file.Open()
	if err != nil {
		return false
	}
	defer fileHeader.Close()

	// Читаем первые 512 байт для определения MIME типа
	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		return false
	}

	mimeType := http.DetectContentType(buffer)
	return allowedTypes[mimeType]
}
