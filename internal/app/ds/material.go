package ds

type Material struct {
	MaterialID          uint   `gorm:"primaryKey;column:material_id"`
	MaterialName        string `gorm:"type:varchar(100);not null"`
	MaterialImageURL    string `gorm:"type:varchar(255)"`
	MaterialDescription string `gorm:"type:text"`
	Isotopes            string `gorm:"type:varchar(100);not null"`
	SampleRequirements  string `gorm:"type:text;not null"`
	IsDeleted           bool   `gorm:"type:boolean;not null;default:false"`
}
