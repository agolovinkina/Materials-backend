package repository

import (
	"errors"
	"fmt"
	"lr2/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetCurrentDatingRequest(userID int) (*ds.MaterialAnalysisRequest, error) {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("creator_id = ? AND request_status = 'draft'", userID).
		Preload("RequestMaterials").
		Preload("RequestMaterials.Material").
		First(&request).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Создаем новую заявку на датирование
			newRequest := &ds.MaterialAnalysisRequest{
				RequestStatus:  "draft",
				CreatorID:      uint(userID),
				Region:         "Балтийское море",
				ExpeditionDate: time.Now(),
				CarbonAge:      "3450 ± 30 BP",
				CreatedAt:      time.Now(),
			}

			err = r.db.Create(newRequest).Error
			if err != nil {
				return nil, fmt.Errorf("ошибка создания заявки: %w", err)
			}

			// Загружаем созданную заявку с отношениями
			err = r.db.Where("request_id = ?", newRequest.RequestID).
				Preload("RequestMaterials").
				Preload("RequestMaterials.Material").
				First(&request).Error
			if err != nil {
				return nil, fmt.Errorf("ошибка загрузки созданной заявки: %w", err)
			}

			return &request, nil
		}
		return nil, fmt.Errorf("ошибка поиска заявки: %w", err)
	}

	return &request, nil
}

func (r *Repository) GetDatingRequestCount(userID int) int {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("creator_id = ? AND request_status = 'draft'", userID).
		First(&request).Error

	if err != nil || request.RequestID == 0 {
		return 0
	}

	var count int64
	err = r.db.Model(&ds.RequestMaterial{}).
		Where("request_id = ?", request.RequestID).
		Count(&count).Error
	if err != nil {
		return 0
	}

	return int(count)
}

func (r *Repository) AddMaterialToDatingRequest(userID int, materialID uint, comment string, probability int, sampleDescription string, sampleWeight float64) error {
	// Получаем текущую заявку или создаём новую
	request, err := r.GetCurrentDatingRequest(userID)
	if err != nil {
		return err
	}

	// Проверяем, не добавлен ли уже этот материал в заявку
	var count int64
	err = r.db.Model(&ds.RequestMaterial{}).
		Where("request_id = ? AND material_id = ?", request.RequestID, materialID).
		Count(&count).Error
	if err != nil {
		return err
	}

	// Если материал уже добавлен, возвращаем ошибку
	if count > 0 {
		return fmt.Errorf("материал уже добавлен в заявку")
	}

	// Добавляем материал в заявку с переданными значениями
	item := ds.RequestMaterial{
		RequestID:          request.RequestID,
		MaterialID:         materialID,
		Comment:            comment,
		ProbabilityPercent: probability,
		CalendarDate:       "",
		SampleDescription:  sampleDescription,
		SampleWeight:       sampleWeight,
		IsPrimary:          false,
	}
	return r.db.Create(&item).Error
}

func (r *Repository) DeleteDatingRequest(requestID uint) error {
	return r.db.Exec("UPDATE material_analysis_requests SET request_status = 'deleted' WHERE request_id = ?", requestID).Error
}

func (r *Repository) GetDatingRequestByID(requestID uint) (*ds.MaterialAnalysisRequest, error) {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("request_id = ? AND request_status != 'deleted'", requestID).
		Preload("RequestMaterials").
		Preload("RequestMaterials.Material").
		First(&request).Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}
