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

func (r *Repository) GetDatingCartCount(userID int) (int, error) {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("creator_id = ? AND request_status = 'draft'", userID).
		First(&request).Error

	if err != nil || request.RequestID == 0 {
		return 0, nil
	}

	var count int64
	err = r.db.Model(&ds.RequestMaterial{}).
		Where("request_id = ?", request.RequestID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
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

func (r *Repository) GetDatingRequests(status string, startDate, endDate time.Time) ([]ds.MaterialAnalysisRequest, error) {
	var requests []ds.MaterialAnalysisRequest
	query := r.db.Where("request_status != 'deleted'") // УБРАЛИ исключение draft

	if status != "" {
		query = query.Where("request_status = ?", status)
	}
	if !startDate.IsZero() {
		query = query.Where("formed_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("formed_at <= ?", endDate)
	}

	err := query.Preload("Creator").Preload("Moderator").Preload("RequestMaterials.Material").Find(&requests).Error
	return requests, err
}

func (r *Repository) UpdateDatingRequest(requestID uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.MaterialAnalysisRequest{}).Where("request_id = ? AND request_status = 'draft'", requestID).Updates(updates).Error
}

func (r *Repository) FormDatingRequest(requestID uint) error {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("request_id = ? AND request_status = 'draft'", requestID).First(&request).Error
	if err != nil {
		return err
	}

	// Проверяем, что есть материалы в заявке
	var count int64
	err = r.db.Model(&ds.RequestMaterial{}).Where("request_id = ?", requestID).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("нельзя сформировать заявку без материалов")
	}

	now := time.Now()
	return r.db.Model(&ds.MaterialAnalysisRequest{}).Where("request_id = ?", requestID).Updates(map[string]interface{}{
		"request_status": "formed",
		"formed_at":      now,
	}).Error
}

func (r *Repository) ProcessDatingRequest(requestID uint, moderatorID uint, action string, calendarDate string) error {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("request_id = ? AND request_status = 'formed'", requestID).First(&request).Error
	if err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"moderator_id": moderatorID,
	}

	switch action {
	case "complete":
		updates["request_status"] = "completed"
		updates["completed_at"] = now
		// Обновляем календарную дату для всех материалов в заявке
		if calendarDate != "" {
			err = r.db.Model(&ds.RequestMaterial{}).Where("request_id = ?", requestID).Update("calendar_date", calendarDate).Error
			if err != nil {
				return err
			}
		}
	case "reject":
		updates["request_status"] = "rejected"
		updates["completed_at"] = now
	default:
		return fmt.Errorf("недопустимое действие: %s", action)
	}

	return r.db.Model(&ds.MaterialAnalysisRequest{}).Where("request_id = ?", requestID).Updates(updates).Error
}

func (r *Repository) RemoveMaterialFromRequest(requestID uint, materialID uint) error {
	return r.db.Where("request_id = ? AND material_id = ?", requestID, materialID).Delete(&ds.RequestMaterial{}).Error
}

func (r *Repository) UpdateRequestMaterial(requestID uint, materialID uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.RequestMaterial{}).Where("request_id = ? AND material_id = ?", requestID, materialID).Updates(updates).Error
}
