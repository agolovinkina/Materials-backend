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

func (r *Repository) CreateMaterial(material *ds.Material) error {
	return r.db.Create(material).Error
}

func (r *Repository) UpdateMaterial(id uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.Material{}).Where("material_id = ? AND is_deleted = false", id).Updates(updates).Error
}

func (r *Repository) DeleteMaterial(id uint) error {
	return r.db.Model(&ds.Material{}).Where("material_id = ?", id).Update("is_deleted", true).Error
}

func (r *Repository) UpdateMaterialImage(id uint, imageURL string) error {
	return r.db.Model(&ds.Material{}).Where("material_id = ?", id).Update("material_image_url", imageURL).Error
}
