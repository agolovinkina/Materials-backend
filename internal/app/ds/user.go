package ds

import (
	"lr4/internal/app/role"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID      uint      `gorm:"primaryKey;column:user_id"`
	Login       string    `gorm:"type:varchar(150);unique;not null"`
	Password    string    `gorm:"type:varchar(128);not null"`
	IsModerator bool      `gorm:"type:boolean;not null;default:false"`
	Role        role.Role `gorm:"type:integer;not null;default:0"`
	UUID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp"`
}

// BeforeCreate hook for GORM
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.New()
	}
	return nil
}
