package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateDatingRequestRequest struct {
	Region         string `json:"region" binding:"required"`
	ExpeditionDate string `json:"expedition_date" binding:"required"`
	CarbonAgeValue int    `json:"carbon_age_value" binding:"required"`
	CarbonAgeError int    `json:"carbon_age_error" binding:"required"`
}

type ProcessDatingRequestRequest struct {
	Action string `json:"action" binding:"required"`
}

// GetDatingRequests получает список заявок на датирование
// @Summary Get dating requests
// @Description Get list of dating requests with filtering
// @Tags Dating Requests
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Success response with requests"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/requests [get]
func (h *Handler) GetDatingRequests(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("пользователь не авторизован"))
		return
	}

	status := ctx.Query("status")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	// Проверяем роль пользователя
	isAdmin := h.IsAdmin(ctx)

	requests, err := h.Repository.GetDatingRequests(userID, isAdmin, status, startDate, endDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, requests)
}

// GetDatingRequestByID получает заявку по ID
// @Summary Get dating request by ID
// @Description Get dating request details by ID
// @Tags Dating Requests
// @Security BearerAuth
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{} "Success response with request"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Router /dating/requests/{id} [get]
func (h *Handler) GetDatingRequestByID(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("пользователь не авторизован"))
		return
	}

	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	isAdmin := h.IsAdmin(ctx)
	request, err := h.Repository.GetDatingRequestByID(uint(id), userID, isAdmin)
	if err != nil {
		if err.Error() == "access denied" {
			h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("доступ запрещен"))
			return
		}
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	h.successResponse(ctx, request)
}

// UpdateDatingRequest обновляет заявку на датирование
// @Summary Update dating request
// @Description Update dating request data
// @Tags Dating Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body UpdateDatingRequestRequest true "Request update data"
// @Success 200 {object} map[string]interface{} "Request updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/requests/{id} [put]
func (h *Handler) UpdateDatingRequest(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req UpdateDatingRequestRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Region == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("регион обязателен для заполнения"))
		return
	}

	if req.ExpeditionDate == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("дата экспедиции обязательна для заполнения"))
		return
	}

	if req.CarbonAgeValue <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("возраст должен быть положительным числом"))
		return
	}

	if req.CarbonAgeError <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("погрешность должна быть положительным числом"))
		return
	}

	expeditionDate, err := time.Parse("2006-01-02", req.ExpeditionDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный формат даты. Используйте YYYY-MM-DD"))
		return
	}

	updates := map[string]interface{}{
		"region":           req.Region,
		"expedition_date":  expeditionDate,
		"carbon_age_value": req.CarbonAgeValue,
		"carbon_age_error": req.CarbonAgeError,
	}

	err = h.Repository.UpdateDatingRequest(uint(id), updates)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Данные заявки успешно обновлены",
	})
}

// FormDatingRequest формирует заявку
// @Summary Form dating request
// @Description Form dating request for processing
// @Tags Dating Requests
// @Security BearerAuth
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{} "Request formed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/requests/{id}/form [put]
func (h *Handler) FormDatingRequest(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetDatingRequestByID(uint(id), 0, true) // Админский доступ для проверки
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	if request.Region == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нельзя сформировать заявку без указания региона"))
		return
	}

	if request.ExpeditionDate.IsZero() {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нельзя сформировать заявку без указания даты экспедиции"))
		return
	}

	if request.CarbonAgeValue == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нельзя сформировать заявку без указания радиоуглеродного возраста"))
		return
	}

	if len(request.RequestMaterials) == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нельзя сформировать заявку без материалов"))
		return
	}

	err = h.Repository.FormDatingRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Заявка успешно сформирована",
	})
}

// ProcessDatingRequest обрабатывает заявку (завершает/отклоняет)
// @Summary Process dating request
// @Description Process dating request (complete/reject)
// @Tags Dating Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body ProcessDatingRequestRequest true "Process action"
// @Success 200 {object} map[string]interface{} "Request processed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Request not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/requests/{id}/process [put]
func (h *Handler) ProcessDatingRequest(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req ProcessDatingRequestRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Action != "complete" && req.Action != "reject" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверное действие. Допустимые значения: complete, reject"))
		return
	}

	request, err := h.Repository.GetDatingRequestByID(uint(id), 0, true) // Админский доступ для модератора
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	var calendarYear, calendarError int
	var formattedDate string

	if req.Action == "complete" {
		// Формируем строку радиоуглеродного возраста для калькулятора
		carbonAgeStr := fmt.Sprintf("%d ± %d BP", request.CarbonAgeValue, request.CarbonAgeError)

		// Используем калькулятор для расчета календарной даты
		result, err := h.CalcService.CalculateProbability(carbonAgeStr, request.Region)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка расчета календарной даты: %w", err))
			return
		}

		calendarYear = result.Year
		calendarError = result.ErrorRange
		formattedDate = result.FormattedDate
	}

	// Получаем ID модератора из контекста и конвертируем в uint
	moderatorID := uint(h.getCurrentUserIDFromContext(ctx))
	err = h.Repository.ProcessDatingRequest(uint(id), moderatorID, req.Action, calendarYear, calendarError)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	response := gin.H{
		"status":  "success",
		"message": "Заявка успешно обработана",
	}

	if req.Action == "complete" {
		response["calendar_year"] = calendarYear
		response["calendar_error"] = calendarError
		response["formatted_date"] = formattedDate
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteDatingRequest удаляет заявку
// @Summary Delete dating request
// @Description Delete dating request by ID
// @Tags Dating Requests
// @Security BearerAuth
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{} "Request deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/requests/{id} [delete]
func (h *Handler) DeleteDatingRequest(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteDatingRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// RemoveMaterialFromRequest удаляет материал из заявки
// @Summary Remove material from request
// @Description Remove material from dating request
// @Tags Dating Requests
// @Security BearerAuth
// @Param requestId path int true "Request ID"
// @Param materialId path int true "Material ID"
// @Success 200 {object} map[string]interface{} "Material removed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/request-materials/{requestId}/{materialId} [delete]
func (h *Handler) RemoveMaterialFromRequest(ctx *gin.Context) {
	requestIDStr := ctx.Param("requestId")
	materialIDStr := ctx.Param("materialId")

	requestID, err1 := strconv.Atoi(requestIDStr)
	materialID, err2 := strconv.Atoi(materialIDStr)

	if err1 != nil || err2 != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err1)
		return
	}

	err := h.Repository.RemoveMaterialFromRequest(uint(requestID), uint(materialID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// UpdateRequestMaterial обновляет материал в заявке
// @Summary Update request material
// @Description Update material details in dating request
// @Tags Dating Requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param requestId path int true "Request ID"
// @Param materialId path int true "Material ID"
// @Param updates body map[string]interface{} true "Update data"
// @Success 200 {object} map[string]interface{} "Material updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /dating/request-materials/{requestId}/{materialId} [put]
func (h *Handler) UpdateRequestMaterial(ctx *gin.Context) {
	requestIDStr := ctx.Param("requestId")
	materialIDStr := ctx.Param("materialId")

	requestID, err1 := strconv.Atoi(requestIDStr)
	materialID, err2 := strconv.Atoi(materialIDStr)

	if err1 != nil || err2 != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err1)
		return
	}

	var updates map[string]interface{}
	if err := ctx.BindJSON(&updates); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err := h.Repository.UpdateRequestMaterial(uint(requestID), uint(materialID), updates)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
