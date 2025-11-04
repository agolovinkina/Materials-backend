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

func (h *Handler) GetDatingRequests(ctx *gin.Context) {
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

	requests, err := h.Repository.GetDatingRequests(status, startDate, endDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, requests)
}

func (h *Handler) GetDatingRequestByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetDatingRequestByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	h.successResponse(ctx, request)
}

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

func (h *Handler) FormDatingRequest(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetDatingRequestByID(uint(id))
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

	request, err := h.Repository.GetDatingRequestByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	var calendarYear, calendarError int

	if req.Action == "complete" {
		carbonAgeStr := fmt.Sprintf("%d ± %d BP", request.CarbonAgeValue, request.CarbonAgeError)

		result, err := h.CalcService.CalculateProbability(carbonAgeStr, request.Region)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		calendarYear = result.Year
		calendarError = result.ErrorRange
	}

	moderatorID := h.getCurrentModeratorID()
	err = h.Repository.ProcessDatingRequest(uint(id), moderatorID, req.Action, calendarYear, calendarError)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	formattedDate := ""
	if calendarYear != 0 {
		era := "AD"
		year := calendarYear
		if calendarYear < 0 {
			era = "BC"
			year = -calendarYear
		}
		formattedDate = fmt.Sprintf("%d %s ± %d лет", year, era, calendarError)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"calendar_year":  calendarYear,
		"calendar_error": calendarError,
		"formatted_date": formattedDate,
	})
}

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
