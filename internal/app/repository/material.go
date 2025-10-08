package repository

import (
	"fmt"
	"lr2/internal/app/ds"
)

func (r *Repository) GetAllMaterials() ([]ds.Material, error) {
	var materials []ds.Material
	err := r.db.Where("is_deleted = false").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	if len(materials) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return materials, nil
}

func (r *Repository) GetMaterialByID(id uint) (*ds.Material, error) {
	var material ds.Material
	err := r.db.Where("material_id = ? AND is_deleted = false", id).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *Repository) SearchMaterialsByName(name string) ([]ds.Material, error) {
	var materials []ds.Material
	err := r.db.Where("material_name ILIKE ? AND is_deleted = false", "%"+name+"%").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}
