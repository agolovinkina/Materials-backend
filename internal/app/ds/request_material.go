package ds

type RequestMaterial struct {
	RequestID          uint    `gorm:"primaryKey"`
	MaterialID         uint    `gorm:"primaryKey"`
	Comment            string  `gorm:"type:text"`
	ProbabilityPercent int     `gorm:"type:integer"`
	CalendarDate       string  `gorm:"type:varchar(100)"`
	SampleDescription  string  `gorm:"type:text;not null"`
	SampleWeight       float64 `gorm:"type:decimal(8,3);not null"`
	IsPrimary          bool    `gorm:"type:boolean;not null;default:false"`

	// Связи
	Request  MaterialAnalysisRequest `gorm:"foreignKey:RequestID"`
	Material Material                `gorm:"foreignKey:MaterialID"`
}
