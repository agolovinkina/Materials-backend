package ds

type RequestMaterial struct {
	RequestID         uint    `gorm:"primaryKey"`
	MaterialID        uint    `gorm:"primaryKey"`
	Comment           string  `gorm:"type:text"`
	CalendarYear      int     `gorm:"type:integer"`
	CalendarError     int     `gorm:"type:integer"`
	SampleDescription string  `gorm:"type:text;not null"`
	SampleWeight      float64 `gorm:"type:decimal(8,3);not null"`
	IsPrimary         bool    `gorm:"type:boolean;not null;default:false"`

	Request  MaterialAnalysisRequest `gorm:"foreignKey:RequestID"`
	Material Material                `gorm:"foreignKey:MaterialID"`
}
