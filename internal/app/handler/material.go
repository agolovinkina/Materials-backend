package handler

import (
	"net/http"
	"strconv"
	"time"

	"lr2/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllMaterials(ctx *gin.Context) {
	userID := 1 // Захардкоженный ID пользователя

	var materials []ds.Material
	var err error

	searchQuery := ctx.Query("searchmaterialquery")
	if searchQuery == "" {
		materials, err = h.Repository.GetAllMaterials()
	} else {
		materials, err = h.Repository.SearchMaterialsByName(searchQuery)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	datingCount := h.Repository.GetDatingRequestCount(userID)
	currentDating, _ := h.Repository.GetCurrentDatingRequest(userID)

	ctx.HTML(http.StatusOK, "materials.html", gin.H{
		"materials":           materials,
		"searchmaterialquery": searchQuery,
		"DatingCount":         datingCount,
		"currentDatingID":     currentDating.RequestID,
		"time":                time.Now().Format("15:04:05"),
	})
}

func (h *Handler) GetMaterialByID(ctx *gin.Context) {
	userID := 1

	strID := ctx.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	material, err := h.Repository.GetMaterialByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	datingCount := h.Repository.GetDatingRequestCount(userID)
	currentDating, _ := h.Repository.GetCurrentDatingRequest(userID)

	ctx.HTML(http.StatusOK, "material.html", gin.H{
		"material":        material,
		"DatingCount":     datingCount,
		"currentDatingID": currentDating.RequestID,
		"time":            time.Now().Format("15:04:05"),
	})
}
