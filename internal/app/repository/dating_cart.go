package repository

import (
	"errors"
	"fmt"
	"lr4/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) CreateDatingRequest(userID int) (*ds.MaterialAnalysisRequest, error) {
	newRequest := &ds.MaterialAnalysisRequest{
		RequestStatus:  "draft",
		CreatorID:      uint(userID),
		Region:         "",
		ExpeditionDate: time.Time{},
		CarbonAgeValue: 0,
		CarbonAgeError: 0,
		CreatedAt:      time.Now(),
		TotalPrice:     0.0,
	}

	err := r.db.Create(newRequest).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка создания заявки: %w", err)
	}

	fmt.Printf("DEBUG: Created new dating request with ID: %d for user %d\n", newRequest.RequestID, userID)
	return newRequest, nil
}

func (r *Repository) GetOrCreateDatingRequest(userID int) (*ds.MaterialAnalysisRequest, error) {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("creator_id = ? AND request_status = 'draft'", userID).First(&request).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("DEBUG: No draft request found for user %d, creating new one\n", userID)
			return r.CreateDatingRequest(userID)
		}
		return nil, fmt.Errorf("ошибка поиска заявки: %w", err)
	}

	fmt.Printf("DEBUG: Found existing draft request ID: %d for user: %d\n", request.RequestID, userID)
	return &request, nil
}

func (r *Repository) AddMaterialToDatingRequest(userID int, materialID uint, comment string, sampleDescription string, sampleWeight float64) error {
	request, err := r.GetOrCreateDatingRequest(userID)
	if err != nil {
		return fmt.Errorf("ошибка получения заявки: %w", err)
	}

	fmt.Printf("DEBUG: Using request ID: %d for user: %d\n", request.RequestID, userID)

	var material ds.Material
	err = r.db.Where("material_id = ? AND is_deleted = false", materialID).First(&material).Error
	if err != nil {
		return fmt.Errorf("материал не найден: %w", err)
	}

	fmt.Printf("DEBUG: Material found: %s (ID: %d)\n", material.MaterialName, material.MaterialID)

	var count int64
	err = r.db.Model(&ds.RequestMaterial{}).
		Where("request_id = ? AND material_id = ?", request.RequestID, materialID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("ошибка проверки материала: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("материал уже добавлен в заявку")
	}

	item := ds.RequestMaterial{
		RequestID:         request.RequestID,
		MaterialID:        materialID,
		Comment:           comment,
		CalendarYear:      0,
		CalendarError:     0,
		SampleDescription: sampleDescription,
		SampleWeight:      sampleWeight,
		IsPrimary:         false,
	}

	err = r.db.Create(&item).Error
	if err != nil {
		return fmt.Errorf("ошибка добавления материала в заявку: %w", err)
	}

	fmt.Printf("DEBUG: Material %d successfully added to request %d\n", materialID, request.RequestID)
	return nil
}

func (r *Repository) GetDatingCartCount(userID int) (int, error) {
	request, err := r.GetOrCreateDatingRequest(userID)
	if err != nil {
		return 0, err
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

func (r *Repository) GetCurrentDatingRequest(userID int) (*ds.MaterialAnalysisRequestDTO, error) {
	request, err := r.GetOrCreateDatingRequest(userID)
	if err != nil {
		return nil, err
	}

	// Загружаем связанные материалы
	err = r.db.Where("request_id = ?", request.RequestID).
		Preload("Material").
		Find(&request.RequestMaterials).Error
	if err != nil {
		return nil, err
	}

	return r.convertToDTO(request)
}
