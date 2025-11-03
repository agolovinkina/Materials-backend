package repository

import (
	"fmt"
	"lr2/internal/app/ds"
	"time"
)

func (r *Repository) GetCurrentDatingRequest(userID int) (*ds.MaterialAnalysisRequestDTO, error) {
	request, err := r.GetOrCreateDatingRequest(userID)
	if err != nil {
		return nil, err
	}

	err = r.db.Where("request_id = ?", request.RequestID).
		Preload("Material").
		Find(&request.RequestMaterials).Error
	if err != nil {
		return nil, err
	}

	return r.convertToDTO(request)
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

func (r *Repository) DeleteDatingRequest(requestID uint) error {
	return r.db.Exec("UPDATE material_analysis_requests SET request_status = 'deleted' WHERE request_id = ?", requestID).Error
}

func (r *Repository) GetDatingRequestByID(requestID uint) (*ds.MaterialAnalysisRequestDTO, error) {
	var request ds.MaterialAnalysisRequest
	err := r.db.Where("request_id = ? AND request_status != 'deleted'", requestID).
		Preload("RequestMaterials").
		Preload("RequestMaterials.Material").
		First(&request).Error

	if err != nil {
		return nil, err
	}

	return r.convertToDTO(&request)
}

func (r *Repository) GetDatingRequests(status string, startDate, endDate time.Time) ([]ds.MaterialAnalysisRequestDTO, error) {
	var requests []ds.MaterialAnalysisRequest
	query := r.db.Where("request_status != 'deleted'")

	if status != "" {
		query = query.Where("request_status = ?", status)
	}
	if !startDate.IsZero() {
		query = query.Where("formed_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("formed_at <= ?", endDate)
	}

	err := query.Preload("RequestMaterials.Material").Find(&requests).Error
	if err != nil {
		return nil, err
	}

	var result []ds.MaterialAnalysisRequestDTO
	for _, req := range requests {
		dto, err := r.convertToDTO(&req)
		if err != nil {
			return nil, err
		}
		result = append(result, *dto)
	}

	return result, nil
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

func (r *Repository) ProcessDatingRequest(requestID uint, moderatorID uint, action string, calendarYear int, calendarError int) error {
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
		if calendarYear > 0 {
			err = r.db.Model(&ds.RequestMaterial{}).Where("request_id = ?", requestID).Updates(map[string]interface{}{
				"calendar_year":  calendarYear,
				"calendar_error": calendarError,
			}).Error
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

func (r *Repository) convertToDTO(request *ds.MaterialAnalysisRequest) (*ds.MaterialAnalysisRequestDTO, error) {
	dto := &ds.MaterialAnalysisRequestDTO{
		RequestID:      int(request.RequestID),
		RequestStatus:  request.RequestStatus,
		CreatedAt:      request.CreatedAt,
		Region:         request.Region,
		ExpeditionDate: request.ExpeditionDate,
		CarbonAgeValue: request.CarbonAgeValue,
		CarbonAgeError: request.CarbonAgeError,
		TotalPrice:     request.TotalPrice,
		FormedAt:       request.FormedAt,
		CompletedAt:    request.CompletedAt,
	}

	var creator ds.User
	if err := r.db.Select("login").Where("user_id = ?", request.CreatorID).First(&creator).Error; err != nil {
		return nil, fmt.Errorf("ошибка загрузки создателя: %w", err)
	}
	dto.CreatorLogin = creator.Login

	if request.ModeratorID != nil {
		var moderator ds.User
		if err := r.db.Select("login").Where("user_id = ?", *request.ModeratorID).First(&moderator).Error; err != nil {
			return nil, fmt.Errorf("ошибка загрузки модератора: %w", err)
		}
		dto.ModeratorLogin = &moderator.Login
	}

	dto.RequestMaterials = make([]ds.RequestMaterialDTO, len(request.RequestMaterials))
	for i, rm := range request.RequestMaterials {
		formattedDate := ""
		if rm.CalendarYear != 0 {
			era := "AD"
			year := rm.CalendarYear
			if rm.CalendarYear < 0 {
				era = "BC"
				year = -rm.CalendarYear
			}
			formattedDate = fmt.Sprintf("%d %s ± %d лет", year, era, rm.CalendarError)
		}

		dto.RequestMaterials[i] = ds.RequestMaterialDTO{
			MaterialID:        rm.MaterialID,
			MaterialName:      rm.Material.MaterialName,
			MaterialImageURL:  rm.Material.MaterialImageURL,
			Comment:           rm.Comment,
			CalendarYear:      rm.CalendarYear,
			CalendarError:     rm.CalendarError,
			FormattedDate:     formattedDate,
			SampleDescription: rm.SampleDescription,
			SampleWeight:      rm.SampleWeight,
			IsPrimary:         rm.IsPrimary,
		}
	}

	return dto, nil
}
