package ds

import (
	"time"
)

type MaterialAnalysisRequest struct {
	RequestID        uint      `gorm:"primaryKey;column:request_id"`
	RequestStatus    string    `gorm:"type:varchar(20);not null"`
	CreatedAt        time.Time `gorm:"not null;default:current_timestamp"`
	CreatorID        uint      `gorm:"not null"`
	FormedAt         *time.Time
	CompletedAt      *time.Time
	ModeratorID      *uint
	Region           string            `gorm:"type:varchar(200);not null"`
	ExpeditionDate   time.Time         `gorm:"type:date;not null"`
	CarbonAge        string            `gorm:"type:varchar(100)"`
	TotalPrice       float64           `gorm:"type:decimal(10,2)"`
	RequestMaterials []RequestMaterial `gorm:"foreignKey:RequestID"`
	Creator          User              `gorm:"foreignKey:CreatorID"`
	Moderator        User              `gorm:"foreignKey:ModeratorID"`
}
