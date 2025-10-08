package handler

import (
	"net/http"
	"strconv"

	"lr2/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddMaterialToDatingRequest(ctx *gin.Context) {
	userID := 1 // Хардкор ID пользователя

	// Получаем ID материала из формы
	materialIDStr := ctx.PostForm("material_id")
	materialID, err := strconv.Atoi(materialIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Добавляем материал в заявку с пустыми значениями
	err = h.Repository.AddMaterialToDatingRequest(userID, uint(materialID),
		"",  // пустой комментарий
		0,   // вероятность 0%
		"",  // пустое описание образца
		0.0) // вес 0.0

	if err != nil {
		if err.Error() == "материал уже добавлен в заявку" {
			ctx.Redirect(http.StatusFound, "/materials?error=material_already_added")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Перенаправляем обратно на страницу материалов
	ctx.Redirect(http.StatusFound, "/materials")
}

func (h *Handler) GetDatingRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetDatingRequestByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	if len(request.RequestMaterials) == 0 {
		ctx.Redirect(http.StatusFound, "/materials")
		return
	}

	// Создаем карту материалов для быстрого доступа
	materialMap := make(map[uint]*ds.Material)
	for i := range request.RequestMaterials {
		materialMap[request.RequestMaterials[i].MaterialID] = &request.RequestMaterials[i].Material
	}

	ctx.HTML(http.StatusOK, "dating.html", gin.H{
		"dating":      request,
		"materialMap": materialMap,
		"time":        request.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

func (h *Handler) DeleteDatingRequest(ctx *gin.Context) {
	requestIDStr := ctx.PostForm("request_id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteDatingRequest(uint(requestID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/materials")
}
