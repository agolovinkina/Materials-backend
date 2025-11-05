package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"lr4/internal/app/ds"

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
// @Success 204 "Material updated successfully"
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

// DeleteMaterial удаляет материал
// @Summary Delete material
// @Description Delete material by ID (admin only)
// @Tags Materials
// @Security BearerAuth
// @Param id path int true "Material ID"
// @Success 204 "Material deleted successfully"
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

	err = h.Repository.DeleteMaterial(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UploadMaterialImage загружает изображение для материала
// @Summary Upload material image
// @Description Upload image for material (admin only)
// @Tags Materials
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Material ID"
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]interface{} "Image uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /materials/{id}/image [post]
func (h *Handler) UploadMaterialImage(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("ошибка получения файла: %w", err))
		return
	}

	if !isImageFile(file) {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("файл должен быть изображением (JPEG, PNG, GIF)"))
		return
	}

	uploadPath := "./uploads/materials/"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("material_%d_%d%s", id, time.Now().Unix(), fileExt)
	filePath := filepath.Join(uploadPath, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	imageURL := fmt.Sprintf("/static/uploads/materials/%s", fileName)

	err = h.Repository.UpdateMaterialImage(uint(id), imageURL)
	if err != nil {
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

// TestFileUpload тестовый endpoint для загрузки файлов
// @Summary Test file upload
// @Description Test endpoint for file upload debugging
// @Tags Debug
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Test file"
// @Success 200 {object} map[string]interface{} "Test successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /test-upload [post]
func (h *Handler) TestFileUpload(ctx *gin.Context) {
	fmt.Println("DEBUG: TestFileUpload called")

	form, err := ctx.MultipartForm()
	if err != nil {
		fmt.Printf("DEBUG: MultipartForm error: %v\n", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	files := form.File["image"]
	fmt.Printf("DEBUG: Files in form: %d\n", len(files))

	if len(files) == 0 {
		fmt.Printf("DEBUG: All form fields: %+v\n", form.Value)
		ctx.JSON(400, gin.H{"error": "No files received", "form_fields": form.Value})
		return
	}

	for i, file := range files {
		fmt.Printf("DEBUG: File %d: %s, Size: %d\n", i, file.Filename, file.Size)
	}

	ctx.JSON(200, gin.H{
		"message":     "test successful",
		"files_count": len(files),
		"files":       files,
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

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		return false
	}

	mimeType := http.DetectContentType(buffer)
	return allowedTypes[mimeType]
}
