package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateDatingRequestRequest struct {
	Region         string `json:"region"`
	ExpeditionDate string `json:"expedition_date"` // ИЗМЕНИЛИ НА string
	CarbonAge      string `json:"carbon_age"`
}

type ProcessDatingRequestRequest struct {
	Action string `json:"action" binding:"required"` // "complete" или "reject"
}

// GetDatingRequests - метод 8
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

// GetDatingRequestByID - метод 9
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

// UpdateDatingRequest - метод 10
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

	updates := make(map[string]interface{})
	if req.Region != "" {
		updates["region"] = req.Region
	}
	if req.ExpeditionDate != "" {
		// Парсим дату из строки
		expeditionDate, err := time.Parse("2006-01-02", req.ExpeditionDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("неверный формат даты. Используйте YYYY-MM-DD"))
			return
		}
		updates["expedition_date"] = expeditionDate
	}
	if req.CarbonAge != "" {
		updates["carbon_age"] = req.CarbonAge
	}

	err = h.Repository.UpdateDatingRequest(uint(id), updates)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// FormDatingRequest - метод 11
func (h *Handler) FormDatingRequest(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.FormDatingRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// ProcessDatingRequest - метод 12
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

	// Получаем заявку для расчета календарной даты
	request, err := h.Repository.GetDatingRequestByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	var calendarDate string
	if req.Action == "complete" {
		// Рассчитываем календарную дату на основе радиоуглеродного анализа
		calendarDate, err = h.CalcService.CalculateCalendarDate(request.CarbonAge, request.Region)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	moderatorID := h.getCurrentModeratorID()
	err = h.Repository.ProcessDatingRequest(uint(id), moderatorID, req.Action, calendarDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"calendar_date": calendarDate,
	})
}

// DeleteDatingRequest - метод 13
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

// RemoveMaterialFromRequest - метод 14
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

// UpdateRequestMaterial - метод 15
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
